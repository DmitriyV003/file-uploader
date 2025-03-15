FROM golang:1.23.4-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
COPY vendor/ vendor/

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .

EXPOSE ${PORT}

CMD ["./main"]
