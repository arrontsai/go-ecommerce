package main

import (
	// 必要套件導入...
	"google.golang.org/grpc"
	"log"
	"net"

	pb "go_project/services/order/proto"
)

func main() {
	// 初始化PostgreSQL
	pgClient, err := database.NewPostgresClient(cfg.PostgresConn)
	if err != nil {
		log.Fatal("無法連接PostgreSQL:", err)
	}

	// 初始化Kafka消費者
	kafkaConsumer, err := messaging.NewMessageBrokerFactory(cfg, zap.L()).CreateMessageBroker(messaging.KafkaBroker)
	if err != nil {
		log.Fatal("無法初始化Kafka消費者:", err)
	}

	// 訂閱Kafka主題
	go kafkaConsumer.Subscribe("order-events", processOrder)

	// 啟動gRPC伺服器
	lis, _ := net.Listen("tcp", ":50051")
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &orderServer{})
	s.Serve(lis)
}
