# Build
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /stockflow ./cmd/api

# Run
FROM alpine:latest
RUN adduser -D stockflow
USER stockflow
COPY --from=builder /stockflow /stockflow
EXPOSE 8080
CMD ["/stockflow"]