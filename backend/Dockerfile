# Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /backend ./cmd/api

# Run
FROM alpine:latest
RUN adduser -D backend
USER stockflow
COPY --from=builder /backend /backend
EXPOSE 8080
CMD ["/backend"]