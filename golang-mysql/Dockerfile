# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder
# Uncomment for Shanghai timezone
# RUN apk add --no-cache tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo Asia/Shanghai > /etc/timezone
ENV CGO_ENABLED=0
WORKDIR /app

# Use Go proxy for faster builds in CN
ENV GOPROXY=https://mirrors.cloud.tencent.com/go/,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /out/server ./cmd/server

FROM alpine:3.20 AS runner
# RUN apk add --no-cache tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo Asia/Shanghai > /etc/timezone
WORKDIR /app

# Static assets
COPY web ./web
COPY --from=builder /out/server ./server

# Copy environment file
COPY .env ./.env

# Port set by env PORT (default 80)
ENV PORT=8080

CMD ["/app/server"]