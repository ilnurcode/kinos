FROM golang:1.26.1-alpine3.23 AS builder

WORKDIR /app

COPY proto ./proto
COPY services/user-service ./services/user-service

RUN go work init ./proto ./services/user-service \
    && cd /app/services/user-service \
    && go build -trimpath -o /out/user_service ./cmd

FROM debian:bookworm-slim

RUN useradd --system --create-home --home-dir /app --shell /usr/sbin/nologin appuser

WORKDIR /app

COPY --from=builder /out/user_service ./user_service
COPY --from=builder /app/services/user-service/migrations ./migrations

USER appuser

EXPOSE 8081

CMD ["./user_service"]
