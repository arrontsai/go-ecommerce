.PHONY: all build clean run test proto docker-build docker-up docker-down kafka-topics kafka-consumer

# Default target
all: build

# Build all services
build:
	@echo "Building all services..."
	@cd services/auth && go build -o ../../bin/auth-service ./cmd/main.go
	@cd services/product && go build -o ../../bin/product-service ./cmd/main.go
	@cd services/order && go build -o ../../bin/order-service ./cmd/main.go
	@cd services/cart && go build -o ../../bin/cart-service ./cmd/main.go
	@cd services/payment && go build -o ../../bin/payment-service ./cmd/main.go
	@echo "Build completed!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "Clean completed!"

# Run all services locally (not in Docker)
run: build
	@echo "Starting all services..."
	@./bin/auth-service &
	@./bin/product-service &
	@./bin/order-service &
	@./bin/cart-service &
	@./bin/payment-service &
	@echo "All services started!"

# Run tests for all services
test:
	@echo "Running tests for all services..."
	@cd services/auth && go test ./...
	@cd services/product && go test ./...
	@cd services/order && go test ./...
	@cd services/cart && go test ./...
	@cd services/payment && go test ./...
	@echo "All tests completed!"

# Generate protobuf files
proto:
	@echo "Generating protobuf files..."
	@protoc --go_out=. --go-grpc_out=. ./proto/*.proto
	@echo "Protobuf files generated!"

# Build Docker images for all services
docker-build:
	@echo "Building Docker images..."
	@docker-compose build
	@echo "Docker images built!"

# Start all services in Docker
docker-up:
	@echo "Starting all services in Docker..."
	@docker-compose up -d
	@echo "All services started in Docker!"

# Stop all services in Docker
docker-down:
	@echo "Stopping all services in Docker..."
	@docker-compose down
	@echo "All services stopped in Docker!"

# Create Kafka topics
kafka-topics:
	@echo "Creating Kafka topics..."
	@docker exec -it kafka kafka-topics --create --topic user-events --bootstrap-server kafka:9092 --partitions 1 --replication-factor 1
	@docker exec -it kafka kafka-topics --create --topic product-events --bootstrap-server kafka:9092 --partitions 1 --replication-factor 1
	@docker exec -it kafka kafka-topics --create --topic order-events --bootstrap-server kafka:9092 --partitions 1 --replication-factor 1
	@docker exec -it kafka kafka-topics --create --topic cart-events --bootstrap-server kafka:9092 --partitions 1 --replication-factor 1
	@docker exec -it kafka kafka-topics --create --topic payment-events --bootstrap-server kafka:9092 --partitions 1 --replication-factor 1
	@echo "Kafka topics created!"

# Run Kafka consumer for debugging
kafka-consumer:
	@echo "Starting Kafka consumer for debugging..."
	@docker exec -it kafka kafka-console-consumer --bootstrap-server kafka:9092 --topic user-events --from-beginning
