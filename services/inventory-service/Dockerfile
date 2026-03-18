FROM golang:1.26.1-alpine3.23 AS builder
WORKDIR /app
COPY . .
WORKDIR /app/services/inventory-service
RUN go mod tidy
RUN go build -o inventory_service ./cmd

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/services/inventory-service/inventory_service .
COPY --from=builder /app/services/inventory-service/migrations ./migrations
EXPOSE 8083
CMD ["./inventory_service"]
