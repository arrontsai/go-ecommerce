# 學習日誌

## 專案概述
本專案是一個基於 Go 語言的微服務電子商務平台，採用現代化的微服務架構設計，實現了產品管理、購物車、訂單處理、支付等核心功能。

## 學習進度

### 2025-03-19
- 初始化專案結構
- 建立基本的服務架構：認證服務、產品服務、訂單服務、購物車服務、支付服務
- 設定共享套件：配置管理、日誌記錄、資料庫連接
- 建立 Docker 環境配置
- 建立 MongoDB 初始化腳本
- 將消息佇列系統從 RabbitMQ 遷移到 Kafka
- 移除所有 RabbitMQ 相關代碼和配置
- 創建前端界面，使用 React 實現簡單的電子商務展示和交互功能

## 微服務間的消息傳遞遷移
將 RabbitMQ 遷移到 Kafka 的主要工作包括：
1. 建立 Kafka 客戶端實現
2. 移除 RabbitMQ 客戶端代碼
3. 更新消息代理工廠類，使其只支持 Kafka
4. 更新配置文件，移除 RabbitMQ 相關配置
5. 更新 Docker Compose 文件，移除 RabbitMQ 服務，確保所有服務使用 Kafka
6. 更新專案依賴，移除 RabbitMQ 相關庫

## 待完成項目
- 實現各微服務的核心功能
- 建立服務間通訊（gRPC 和 Kafka 消息佇列）
- 實現 API 閘道
- 完成前端界面的剩餘功能
- 部署到 Kubernetes 環境

## 學習資源
- Go 官方文檔：https://golang.org/doc/
- gRPC 文檔：https://grpc.io/docs/
- Docker 和 Docker Compose 文檔：https://docs.docker.com/
- MongoDB 文檔：https://docs.mongodb.com/
- PostgreSQL 文檔：https://www.postgresql.org/docs/
- Kafka 文檔：https://kafka.apache.org/documentation/
- React 文檔：https://reactjs.org/docs/

## 問題與解決方案
### 消息佇列遷移
問題：在將 RabbitMQ 遷移到 Kafka 的過程中，需要處理不同的消息傳遞模型和 API。
解決方案：設計了一個抽象的工廠模式，使服務代碼可以透過統一的介面與消息系統交互，從而實現了平滑遷移。

### 同時支持多種消息傳遞系統
問題：在遷移期間，需要同時支持 RabbitMQ 和 Kafka 以確保系統穩定性。
解決方案：使用了工廠模式來創建不同的消息客戶端，並使用偏好設置來決定默認使用哪一個。完成遷移後，移除了 RabbitMQ 相關代碼。
