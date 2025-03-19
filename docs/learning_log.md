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

### 2025-03-20
- 實現購物車服務 (Cart Service) 的 REST API
- 實現購物車儲存庫 (CartRepository) 介面和 MongoDB 實現
- 定義購物車資料模型 (Cart 和 CartItem)
- 實現訂單服務 (Order Service) 的 gRPC 接口
- 定義訂單 protobuf 文件並生成 Go 代碼
- 實現訂單資料模型 (Order 和 OrderItem)
- 實現 Kafka 事件處理，使購物車結帳事件能夠觸發訂單創建

## 微服務間的消息傳遞遷移
將 RabbitMQ 遷移到 Kafka 的主要工作包括：
1. 建立 Kafka 客戶端實現
2. 移除 RabbitMQ 客戶端代碼
3. 更新消息代理工廠類，使其只支持 Kafka
4. 更新配置文件，移除 RabbitMQ 相關配置
5. 更新 Docker Compose 文件，移除 RabbitMQ 服務，確保所有服務使用 Kafka
6. 更新專案依賴，移除 RabbitMQ 相關庫

## 待完成項目
- 實現認證服務 (Auth Service) 的 JWT 認證
- 實現產品服務 (Product Service) 的 REST API
- 實現支付服務 (Payment Service) 的支付處理
- 建立 API 閘道整合各微服務
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

### 微服務間的事件驅動通訊
問題：購物車結帳後，需要將結帳事件傳遞給訂單服務，以便創建新訂單。
解決方案：使用 Kafka 作為事件總線，購物車服務在結帳時發布事件，訂單服務訂閱相關主題並處理這些事件，實現了鬆耦合的服務間通訊。

### Protobuf 和 gRPC 設置
問題：設置 gRPC 服務需要正確配置 protobuf 文件和生成代碼。
解決方案：定義了清晰的 protobuf 文件，並使用 protoc 工具生成 Go 代碼。確保了 go_package 選項正確設置，以便生成的代碼能夠被正確導入。
