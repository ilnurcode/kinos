FROM golang:1.26.1-alpine3.23 AS builder
WORKDIR /app
COPY proto ./proto
COPY services/inventory-service ./services/inventory-service
RUN go work init ./proto ./services/inventory-service \
    && cd /app/services/inventory-service \
    && go build -trimpath -o /out/inventory_service ./cmd

FROM debian:bookworm-slim
RUN useradd --system --create-home --home-dir /app --shell /usr/sbin/nologin appuser
WORKDIR /app
COPY --from=builder /out/inventory_service ./inventory_service
COPY --from=builder /app/services/inventory-service/migrations ./migrations
USER appuser
EXPOSE 8083
CMD ["./inventory_service"]
