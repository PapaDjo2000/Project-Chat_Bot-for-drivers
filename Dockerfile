FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/app/main.go

FROM alpine:latest

RUN apk add --no-cache bash

COPY --from=builder /app/main /app/main

WORKDIR /app

EXPOSE 8080

CMD ["./main"]