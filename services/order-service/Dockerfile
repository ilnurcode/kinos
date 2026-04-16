FROM golang:1.26.1-alpine3.23 AS builder

WORKDIR /app

COPY proto ./proto
COPY services/order-service ./services/order-service

RUN go work init ./proto ./services/order-service \
    && cd /app/services/order-service \
    && go build -trimpath -o /out/order_service ./cmd/main.go

FROM debian:bookworm-slim

RUN useradd --system --create-home --home-dir /app --shell /usr/sbin/nologin appuser

WORKDIR /app

COPY --from=builder /out/order_service ./order_service
COPY --from=builder /app/services/order-service/migrations ./migrations

USER appuser

EXPOSE 8085

CMD ["./order_service"]
