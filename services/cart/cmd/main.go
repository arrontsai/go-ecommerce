package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ecommerce/pkg/config"
	"github.com/yourusername/ecommerce/pkg/database"
	"github.com/yourusername/ecommerce/pkg/logger"
	"github.com/yourusername/ecommerce/pkg/middleware"
	"github.com/yourusername/ecommerce/services/cart/handler"
	"github.com/yourusername/ecommerce/services/cart/repository"
	"github.com/yourusername/ecommerce/services/cart/service"
	"go.uber.org/zap"
)

func main() {
	// 初始化配置
	cfg := config.LoadConfig()

	// 初始化MongoDB
	mongoClient, err := database.NewMongoClient(cfg.MongoURI)
	if err != nil {
		log.Fatal("無法連接MongoDB:", err)
	}

	// 初始化購物車倉儲
	cartRepo := repository.NewCartRepo(mongoClient, cfg.MongoDB)

	// 初始化Kafka生產者
	kafkaProducer, err := messaging.NewMessageBrokerFactory(cfg, zap.L()).CreateMessageBroker(messaging.KafkaBroker)
	if err != nil {
		log.Fatal("無法初始化Kafka生產者:", err)
	}

	// 啟動HTTP伺服器
	r := gin.Default()
	r.POST("/cart", handler.AddToCart(cartRepo, kafkaProducer))
	r.Run(":8083")
}
