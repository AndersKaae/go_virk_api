# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy the local dependency
COPY go_virk_updater /go_virk_updater

# Copy go mod files
COPY go_virk_api/go.mod go_virk_api/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY go_virk_api/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api cmd/api/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/api .

# Expose port
EXPOSE 8444

# Run the application
CMD ["./api"]
