FROM golang:1.24.5 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o order-service ./services/order-service/cmd

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/order-service .

EXPOSE 8081

ENTRYPOINT ["./order-service"]