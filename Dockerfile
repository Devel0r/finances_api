FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV DB_URL=postgres://postgres:postgres@postgres:5432/finance_api?sslmode=disable

RUN CGO_ENABLED=0 GOOS=linux go build -o finance_api cmd/finance_api/main.go 

CMD ["./finance_api"]