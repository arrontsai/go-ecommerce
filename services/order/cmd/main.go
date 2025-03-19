package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/arrontsai/ecommerce/pkg/config"
	"github.com/arrontsai/ecommerce/pkg/logger"
	"github.com/arrontsai/ecommerce/pkg/messaging"
	"github.com/arrontsai/ecommerce/services/order/model"
	"github.com/arrontsai/ecommerce/services/order/proto/pb"
	"github.com/arrontsai/ecommerce/services/order/repository"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// orderServer 實現訂單服務的gRPC接口
type orderServer struct {
	pb.UnimplementedOrderServiceServer
	db *sqlx.DB
}

func main() {
	// 初始化配置
	cfg, err := config.LoadConfig("order-service")
	if err != nil {
		log.Fatal("無法載入配置:", err)
	}

	// 初始化日誌記錄器
	appLogger, err := logger.NewLogger("order-service", cfg.LogLevel)
	if err != nil {
		log.Fatal("無法初始化日誌記錄器:", err)
	}

	// 構建 PostgreSQL 連接字符串
	pgConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// 初始化PostgreSQL
	pgClient, err := sqlx.Connect("postgres", pgConnStr)
	if err != nil {
		appLogger.Fatal("無法連接PostgreSQL:", zap.Error(err))
	}
	defer pgClient.Close()

	// 初始化Kafka消費者
	kafkaConsumer, err := messaging.NewKafkaClient(cfg.KafkaBrokers, appLogger.Logger)
	if err != nil {
		appLogger.Fatal("無法初始化Kafka消費者:", zap.Error(err))
	}

	// 創建訂單服務
	server := &orderServer{db: pgClient}

	// 訂閱Kafka主題
	go subscribeToCartEvents(kafkaConsumer, server)

	// 啟動gRPC伺服器
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		appLogger.Fatal("無法監聽端口:", zap.Error(err))
	}
	
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, server)
	
	appLogger.Info("訂單服務啟動於 :50051")
	if err := s.Serve(lis); err != nil {
		appLogger.Fatal("無法啟動gRPC伺服器:", zap.Error(err))
	}
}

// CreateOrder 實現創建訂單的gRPC方法
func (s *orderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	// 生成訂單ID
	orderID := fmt.Sprintf("%x", os.Getpid())

	// 開始資料庫事務
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("無法開始事務: %v", err)
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
		return nil, fmt.Errorf("無法創建訂單: %v", err)
	}
	
	// 插入訂單項目
	for _, item := range req.Items {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO order_items (order_id, product_id, quantity, price) 
			 VALUES ($1, $2, $3, $4)`,
			orderID, item.ProductId, item.Quantity, item.Price)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("無法創建訂單項目: %v", err)
		}
	}
	
	// 提交事務
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("無法提交事務: %v", err)
	}
	
	return &pb.OrderResponse{
		OrderId: orderID,
		Status:  "CREATED",
	}, nil
}

// GetOrder 實現獲取訂單的gRPC方法
func (s *orderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderDetailResponse, error) {
	// 從數據庫獲取訂單
	var order struct {
		ID          string  `db:"order_id"`
		UserID      string  `db:"user_id"`
		TotalAmount float64 `db:"total_price"`
		Status      string  `db:"status"`
	}

	err := s.db.Get(&order, "SELECT order_id, user_id, total_price, status FROM orders WHERE order_id = $1", req.OrderId)
	if err != nil {
		return nil, fmt.Errorf("獲取訂單失敗: %w", err)
	}

	// 獲取訂單項目
	var items []struct {
		ProductID   string  `db:"product_id"`
		ProductName string  `db:"product_name"`
		Quantity    int32   `db:"quantity"`
		UnitPrice   float64 `db:"unit_price"`
	}

	err = s.db.Select(&items, "SELECT product_id, product_name, quantity, unit_price FROM order_items WHERE order_id = $1", req.OrderId)
	if err != nil {
		return nil, fmt.Errorf("獲取訂單項目失敗: %w", err)
	}

	// 轉換為gRPC響應格式
	pbItems := make([]*pb.OrderItem, 0, len(items))
	for _, item := range items {
		pbItems = append(pbItems, &pb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
			Price:     float32(item.UnitPrice),
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
	// 更新訂單狀態
	_, err := s.db.Exec("UPDATE orders SET status = $1 WHERE order_id = $2", req.Status, req.OrderId)
	if err != nil {
		return nil, fmt.Errorf("更新訂單狀態失敗: %w", err)
	}

	return &pb.OrderResponse{
		OrderId: req.OrderId,
		Status:  req.Status,
	}, nil
}

// subscribeToCartEvents 訂閱購物車事件
func subscribeToCartEvents(consumer *messaging.KafkaClient, server *orderServer) {
	// 處理消息的回調函數
	handler := func(msg []byte) error {
		var event struct {
			EventType string          `json:"event_type"`
			UserID    string          `json:"user_id"`
			Items     []model.OrderItem `json:"items"`
		}

		err := json.Unmarshal(msg, &event)
		if err != nil {
			log.Printf("解析消息失敗: %v", err)
			return err
		}

		// 根據事件類型處理不同的邏輯
		if event.EventType == "CHECKOUT" {
			// 創建訂單
			err := repository.NewOrderRepo(server.db).CreateOrder(event.UserID, event.Items)
			if err != nil {
				log.Printf("創建訂單失敗: %v", err)
				return err
			}
			log.Printf("已為用戶 %s 創建訂單", event.UserID)
		}

		return nil
	}

	// 啟動消費
	ctx := context.Background()
	err := consumer.ConsumeMessages(ctx, "cart_events", "order-service", handler)
	if err != nil {
		log.Fatalf("無法消費消息: %v", err)
	}
}
