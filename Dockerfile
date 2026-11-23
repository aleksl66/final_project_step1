# Этап 1: Сборка
FROM golang:1.26.1 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

# Этап 2: Финальный образ
FROM ubuntu:latest
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/server .
COPY web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db

EXPOSE 7540
CMD ["./server"]