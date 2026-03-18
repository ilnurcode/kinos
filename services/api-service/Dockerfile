FROM golang:1.26.1-alpine3.23 AS builder

WORKDIR /app

COPY . .

WORKDIR /app/services/api-service

RUN go mod tidy
RUN go build -o api_service ./cmd

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/services/api-service/api_service .
COPY --from=builder /app/services/api-service/templates ./templates
COPY --from=builder /app/services/api-service/static ./static

EXPOSE 8080

CMD ["./api_service"]