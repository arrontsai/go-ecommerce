package main

import (
	"context"
	"encoding/json"
	"log"
	"net"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/Shopify/sarama"
	"github.com/jmoiron/sqlx"

	"github.com/arrontsai/ecommerce/pkg/config"
	"github.com/arrontsai/ecommerce/pkg/database"
	"github.com/arrontsai/ecommerce/pkg/messaging"
	"github.com/arrontsai/ecommerce/pkg/models"
	pb "github.com/arrontsai/ecommerce/services/order/proto/pb"
)

// orderServer 實現訂單服務的gRPC接口
type orderServer struct {
	pb.UnimplementedOrderServiceServer
	db *sqlx.DB
}

func main() {
	// 初始化配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("無法載入配置:", err)
	}

	// 初始化PostgreSQL
	pgClient, err := database.NewPostgresClient(cfg.PostgreSQL.URI)
	if err != nil {
		log.Fatal("無法連接PostgreSQL:", err)
	}
	defer pgClient.Close()

	// 初始化Kafka消費者
	kafkaFactory := messaging.NewMessageBrokerFactory(cfg)
	kafkaConsumer, err := kafkaFactory.CreateKafkaConsumer("order-group")
	if err != nil {
		log.Fatal("無法初始化Kafka消費者:", err)
	}

	// 創建訂單服務
	server := &orderServer{db: pgClient}

	// 訂閱Kafka主題
	go subscribeToCartEvents(kafkaConsumer, server)

	// 啟動gRPC伺服器
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("無法監聽端口: %v", err)
	}
	
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, server)
	
	log.Println("訂單服務啟動於 :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("無法啟動gRPC伺服器: %v", err)
	}
}

// CreateOrder 實現創建訂單的gRPC方法
func (s *orderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	// 生成訂單ID
	orderID := uuid.New().String()
	
	// 開始資料庫事務
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "無法開始事務: %v", err)
	}
	
	// 計算總金額
	var totalAmount float64
	for _, item := range req.Items {
		totalAmount += float64(item.Price) * float64(item.Quantity)
	}
	
	// 插入訂單主記錄
	_, err = tx.ExecContext(ctx,
		`INSERT INTO orders (id, user_id, total_amount, status, created_at) 
		 VALUES ($1, $2, $3, $4, NOW())`,
		orderID, req.UserId, totalAmount, "PENDING")
	if err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "無法創建訂單: %v", err)
	}
	
	// 插入訂單項目
	for _, item := range req.Items {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO order_items (order_id, product_id, quantity, price) 
			 VALUES ($1, $2, $3, $4)`,
			orderID, item.ProductId, item.Quantity, item.Price)
		if err != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "無法創建訂單項目: %v", err)
		}
	}
	
	// 提交事務
	if err = tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "無法提交事務: %v", err)
	}
	
	return &pb.OrderResponse{
		OrderId: orderID,
		Status:  "CREATED",
	}, nil
}

// GetOrder 實現獲取訂單的gRPC方法
func (s *orderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderDetailResponse, error) {
	var order models.Order
	err := s.db.GetContext(ctx, &order,
		`SELECT id, user_id, total_amount, status, created_at 
		 FROM orders WHERE id = $1`, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "訂單不存在: %v", err)
	}
	
	// 獲取訂單項目
	var items []models.OrderItem
	err = s.db.SelectContext(ctx, &items,
		`SELECT product_id, quantity, price 
		 FROM order_items WHERE order_id = $1`, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "無法獲取訂單項目: %v", err)
	}
	
	// 轉換為gRPC響應
	var pbItems []*pb.OrderItem
	for _, item := range items {
		pbItems = append(pbItems, &pb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     float32(item.Price),
		})
	}
	
	return &pb.OrderDetailResponse{
		OrderId:     order.ID,
		UserId:      order.UserID,
		TotalAmount: float32(order.TotalAmount),
		Status:      order.Status,
		Items:       pbItems,
	}, nil
}

// UpdateOrderStatus 實現更新訂單狀態的gRPC方法
func (s *orderServer) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.OrderResponse, error) {
	result, err := s.db.ExecContext(ctx,
		`UPDATE orders SET status = $1 WHERE id = $2`,
		req.Status, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "無法更新訂單狀態: %v", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "無法獲取影響的行數: %v", err)
	}
	
	if rowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "訂單不存在")
	}
	
	return &pb.OrderResponse{
		OrderId: req.OrderId,
		Status:  req.Status,
	}, nil
}

// subscribeToCartEvents 訂閱購物車事件
func subscribeToCartEvents(consumer messaging.KafkaConsumer, server *orderServer) {
	// 處理消息的回調函數
	handler := func(msg *sarama.ConsumerMessage) error {
		var event struct {
			EventType string          `json:"event_type"`
			UserID    string          `json:"user_id"`
			CartID    string          `json:"cart_id"`
			Items     []models.CartItem `json:"items"`
		}
		
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("無法解析事件: %v", err)
			return err
		}
		
		// 只處理結帳事件
		if event.EventType != "CHECKOUT" {
			return nil
		}
		
		// 將購物車項目轉換為訂單項目
		var orderItems []*pb.OrderItem
		for _, item := range event.Items {
			// 這裡應該從產品服務獲取價格，但為了簡化，我們假設價格為10.0
			orderItems = append(orderItems, &pb.OrderItem{
				ProductId: item.ProductID,
				Quantity:  int32(item.Quantity),
				Price:     10.0, // 假設價格
			})
		}
		
		// 創建訂單
		_, err := server.CreateOrder(context.Background(), &pb.CreateOrderRequest{
			UserId: event.UserID,
			Items:  orderItems,
		})
		
		if err != nil {
			log.Printf("無法創建訂單: %v", err)
			return err
		}
		
		log.Printf("已為用戶 %s 創建訂單", event.UserID)
		return nil
	}
	
	// 開始消費消息
	err := consumer.Consume("cart-events", handler)
	if err != nil {
		log.Fatalf("無法消費消息: %v", err)
	}
}
