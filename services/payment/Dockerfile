FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payment-service ./services/payment/cmd

FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/payment-service .

# Run the service
CMD ["./payment-service"]
