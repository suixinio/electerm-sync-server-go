FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go module files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/electerm-sync-server ./src

# Create a minimal runtime image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/bin/electerm-sync-server /app/electerm-sync-server
COPY --from=builder /app/sample.env /app/.env

# Create a directory for storing data if FILE_STORE_PATH is not explicitly set
RUN mkdir -p /app/data

# Set environment variables with defaults that can be overridden
ENV PORT=7837
ENV HOST=0.0.0.0
ENV FILE_STORE_PATH=/app/data

# Expose the port the service runs on
EXPOSE 7837

# Run the binary
CMD ["/app/electerm-sync-server"]