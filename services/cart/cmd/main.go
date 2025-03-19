package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/arrontsai/ecommerce/pkg/config"
	"github.com/arrontsai/ecommerce/pkg/database"
	"github.com/arrontsai/ecommerce/pkg/messaging"
	"github.com/arrontsai/ecommerce/pkg/models"
	"github.com/arrontsai/ecommerce/services/cart/repository"
)

func main() {
	// 初始化配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("無法載入配置:", err)
	}

	// 初始化MongoDB連接
	mongoClient, err := database.NewMongoClient(cfg.MongoDB.URI)
	if err != nil {
		log.Fatal("無法連接MongoDB:", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// 初始化購物車儲存庫
	cartRepo := repository.NewMongoCartRepository(mongoClient)

	// 初始化Kafka生產者
	kafkaFactory := messaging.NewMessageBrokerFactory(cfg)
	kafkaProducer, err := kafkaFactory.CreateKafkaProducer()
	if err != nil {
		log.Fatal("無法初始化Kafka生產者:", err)
	}

	// 設置HTTP路由
	router := setupRouter(cartRepo, kafkaProducer)

	// 啟動HTTP服務器
	log.Println("購物車服務啟動於 :8082")
	if err := router.Run(":8082"); err != nil {
		log.Fatal("HTTP服務器啟動失敗:", err)
	}
}

func setupRouter(repo repository.CartRepository, producer messaging.KafkaProducer) *gin.Engine {
	r := gin.Default()

	// 健康檢查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	// 添加商品到購物車
	r.POST("/cart", addToCartHandler(repo))
	
	// 獲取購物車內容
	r.GET("/cart/:userID", getCartHandler(repo))
	
	// 結帳
	r.POST("/cart/checkout", checkoutHandler(repo, producer))

	return r
}

func addToCartHandler(repo repository.CartRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID    string `json:"user_id" binding:"required"`
			ProductID string `json:"product_id" binding:"required"`
			Quantity  int    `json:"quantity" binding:"required,min=1"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := repo.AddToCart(req.UserID, req.ProductID, req.Quantity); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法添加商品到購物車"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "商品已添加到購物車"})
	}
}

func getCartHandler(repo repository.CartRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")
		cart, err := repo.GetCart(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "購物車不存在"})
			return
		}

		c.JSON(http.StatusOK, cart)
	}
}

func checkoutHandler(repo repository.CartRepository, producer messaging.KafkaProducer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID string `json:"user_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 獲取購物車
		cart, err := repo.GetCart(req.UserID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "購物車不存在"})
			return
		}

		// 發送事件到Kafka
		event := map[string]interface{}{
			"event_type": "CHECKOUT",
			"user_id":    req.UserID,
			"cart_id":    cart.ID,
			"items":      cart.Items,
			"timestamp":  time.Now(),
		}

		if err := producer.Produce("cart-events", event); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "結帳處理失敗"})
			return
		}

		// 清空購物車
		if err := repo.ClearCart(req.UserID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法清空購物車"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   "結帳成功",
			"timestamp": time.Now(),
		})
	}
}
