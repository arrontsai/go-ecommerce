# Go Microservices E-commerce Platform

A modern e-commerce platform built with Go microservices architecture.

## Project Overview

This project is a microservices-based e-commerce platform built with Golang. It follows best practices for microservices architecture and uses modern technologies for service communication, data storage, and deployment.

## Architecture

The application is divided into the following microservices:

1. **Auth Service**: Handles user authentication and authorization
2. **Product Service**: Manages product catalog and inventory
3. **Order Service**: Processes and manages customer orders
4. **Cart Service**: Manages shopping cart functionality
5. **Payment Service**: Handles payment processing

## Technology Stack

- **Backend**: Go (Golang)
- **API**: REST and gRPC
- **Databases**: MongoDB, PostgreSQL
- **Message Broker**: Kafka
- **Containerization**: Docker & Docker Compose
- **Service Discovery**: Consul (planned)
- **API Gateway**: (planned)
- **Monitoring**: Prometheus & Grafana (planned)

## Project Structure

```
go_project/
├── services/             # Individual microservices
│   ├── auth/             # Authentication service
│   ├── product/          # Product catalog service
│   ├── order/            # Order management service
│   ├── cart/             # Shopping cart service
│   └── payment/          # Payment processing service
├── pkg/                  # Shared packages
│   ├── common/           # Common utilities
│   ├── config/           # Configuration management
│   ├── database/         # Database connections
│   ├── logger/           # Logging utilities
│   ├── middleware/       # Shared middleware
│   └── models/           # Shared data models
├── api/                  # API definitions (REST, gRPC)
├── docs/                 # Documentation
├── scripts/              # Utility scripts
└── deployments/          # Deployment configurations
```

## Getting Started

### Prerequisites

- Go 1.20+
- Docker and Docker Compose
- Make (optional)

### Setup and Installation

1. Clone the repository
2. Navigate to the project directory
3. Run the services using Docker Compose:

```bash
docker-compose up -d
```

## Development

Each service can be developed and tested independently. Refer to the README in each service directory for specific instructions.

## License

MIT
