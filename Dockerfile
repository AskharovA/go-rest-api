# Builder stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s -extldflags "-static"' -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite

RUN addgroup -g 1001 appgroup && \
    adduser -u 1001 -G appgroup -s /bin/sh -D appuser

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/env.example ./

RUN mkdir -p /app/data && chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

CMD ["./main"]