FROM golang:1.26.1-alpine3.23 AS builder

WORKDIR /app

COPY proto ./proto
COPY services/catalog-service ./services/catalog-service

RUN go work init ./proto ./services/catalog-service \
    && cd /app/services/catalog-service \
    && go build -trimpath -o /out/catalog_service ./cmd

FROM debian:bookworm-slim

RUN useradd --system --create-home --home-dir /app --shell /usr/sbin/nologin appuser

WORKDIR /app

COPY --from=builder /out/catalog_service ./catalog_service
COPY --from=builder /app/services/catalog-service/migrations ./migrations

USER appuser

EXPOSE 8082

CMD ["./catalog_service"]
