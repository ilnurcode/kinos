FROM golang:1.25.7-alpine3.23 AS builder

WORKDIR /app

COPY . .


WORKDIR /app/services/user-service

RUN go mod tidy
RUN go build -o user_service ./cmd

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/services/user-service/user_service .
COPY --from=builder /app/services/user-service/migrations ./migrations

EXPOSE 8081

CMD ["./user_service"]
