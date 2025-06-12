FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o blog-api ./cmd/api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/blog-api .
COPY .env .env

EXPOSE 8088

CMD ["./blog-api"]