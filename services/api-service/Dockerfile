FROM golang:1.26.1-alpine3.23 AS builder

WORKDIR /app

COPY proto ./proto
COPY services/api-service ./services/api-service

RUN go work init ./proto ./services/api-service \
    && cd /app/services/api-service \
    && go build -trimpath -o /out/api_service ./cmd

FROM debian:bookworm-slim

RUN useradd --system --create-home --home-dir /app --shell /usr/sbin/nologin appuser

WORKDIR /app

COPY --from=builder /out/api_service ./api_service

USER appuser

EXPOSE 8080

CMD ["./api_service"]
