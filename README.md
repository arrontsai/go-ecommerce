# Go 微服務電子商務平台

一個使用 Go 微服務架構構建的現代電子商務平台。

## 專案概述

本專案是一個基於微服務架構的電子商務平台，使用 Golang 構建。它遵循微服務架構的最佳實踐，並使用現代技術進行服務通信、數據存儲和部署。

## 架構

應用程序分為以下微服務：

1. **認證服務 (Auth Service)**: 處理用戶認證和授權
2. **產品服務 (Product Service)**: 管理產品目錄和庫存
3. **訂單服務 (Order Service)**: 處理和管理客戶訂單
4. **購物車服務 (Cart Service)**: 管理購物車功能
5. **支付服務 (Payment Service)**: 處理支付流程

## 技術棧

- **後端**: Go (Golang)
- **API**: REST 和 gRPC
- **數據庫**: MongoDB, PostgreSQL
- **消息代理**: Kafka
- **容器化**: Docker & Docker Compose
- **服務發現**: Consul (計劃中)
- **API 網關**: (計劃中)
- **監控**: Prometheus & Grafana (計劃中)

## 專案結構

```
go-ecommerce/
├── services/             # 個別微服務
│   ├── auth/             # 認證服務
│   │   ├── cmd/          # 命令行入口
│   │   ├── handler/      # HTTP 處理器
│   │   ├── repository/   # 數據存儲層
│   │   └── service/      # 業務邏輯層
│   ├── product/          # 產品目錄服務
│   ├── order/            # 訂單管理服務
│   ├── cart/             # 購物車服務
│   └── payment/          # 支付處理服務
├── pkg/                  # 共享套件
│   ├── config/           # 配置管理
│   ├── database/         # 數據庫連接
│   ├── logger/           # 日誌工具
│   ├── messaging/        # 消息傳遞
│   ├── middleware/       # 共享中間件
│   └── models/           # 共享數據模型
├── frontend/             # 前端應用
│   ├── public/           # 靜態資源
│   └── src/              # 源代碼
├── docs/                 # 文檔
└── scripts/              # 工具腳本
```

## 開始使用

### 前提條件

- Go 1.20+
- Docker 和 Docker Compose
- Make (可選)

### 設置和安裝

1. 克隆存儲庫
2. 導航到項目目錄
3. 使用 Docker Compose 運行服務：

```bash
docker-compose up -d
```

## 開發

每個服務都可以獨立開發和測試。有關具體說明，請參閱每個服務目錄中的 README。

## 許可證

MIT
