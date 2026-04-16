FROM golang:1.26.1-alpine3.23 AS builder

WORKDIR /app

COPY proto ./proto
COPY services/cart-service ./services/cart-service

RUN go work init ./proto ./services/cart-service \
    && cd /app/services/cart-service \
    && go build -trimpath -o /out/cart_service ./cmd/main.go

FROM debian:bookworm-slim

RUN useradd --system --create-home --home-dir /app --shell /usr/sbin/nologin appuser

WORKDIR /app

COPY --from=builder /out/cart_service ./cart_service

USER appuser

EXPOSE 8084

CMD ["./cart_service"]
