# Этап сборки
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/finances_api cmd/finances_api/main.go

# Этап выполнения
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /bin/finances_api .

COPY .env.local .env.local

RUN chmod +x finances_api

ENV DB_URL=postgres://postgres:postgres@postgres:5432/finance_api?sslmode=disable

EXPOSE 8080

CMD ["./finances_api"]