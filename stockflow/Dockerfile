# Stage 1: Build binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=1.0.0" \
    -o /app/stockflow-api ./cmd/api

# Stage 2: Minimal runtime
FROM alpine:3.19 AS runner

RUN apk --no-cache add ca-certificates tzdata \
    && addgroup -S stockflow \
    && adduser -S -G stockflow -u 10001 stockflow

WORKDIR /app

COPY --from=builder /app/stockflow-api /app/stockflow-api
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

USER 10001:10001

EXPOSE 8080

HEALTHCHECK --interval=15s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/stockflow-api"]