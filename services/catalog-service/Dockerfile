FROM golang:1.25.7-alpine3.23 AS builder

WORKDIR /app

COPY . .


WORKDIR /app/services/catalog-service

RUN go mod tidy
RUN go build -o catalog_service ./cmd

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/services/catalog-service/catalog_service .
COPY --from=builder /app/services/catalog-service/migrations ./migrations

EXPOSE 8082

CMD ["./catalog_service"]
