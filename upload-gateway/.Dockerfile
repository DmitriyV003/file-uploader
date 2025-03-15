FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN apk update && apk add --no-cache git

COPY go.mod go.sum ./
COPY vendor/ vendor/

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main ./cmd

RUN go install github.com/gobuffalo/soda@latest

FROM alpine:latest
WORKDIR /app

RUN apk update && apk add --no-cache ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /app/main .
COPY --from=builder /app/database.yml .
COPY --from=builder /app/migrations migrations/
COPY --from=builder /go/bin/soda /usr/local/bin/soda
COPY wait-for-db.sh .
RUN chmod +x wait-for-db.sh

EXPOSE ${PORT}

CMD ["./wait-for-db.sh", "mysql", "./main"]
