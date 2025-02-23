# 🚀 Stage 1: Build the Go binary
FROM golang:1.21 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app (static binary)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /odyscan ./cmd/main.go

# 🏗️ Stage 2: Create a minimal runtime image
FROM alpine:latest

# Install necessary dependencies (e.g., ClamAV client if needed)
RUN apk --no-cache add ca-certificates clamav

# Set the working directory
WORKDIR /app

# Copy the compiled binary from builder stage
COPY --from=builder /odyscan /odyscan

# Copy configuration files (optional)
COPY config.yaml /app/config.yaml

# Expose necessary ports (adjust if needed)
EXPOSE 8080

# Set environment variables for K3s (if needed)
ENV CONFIG_PATH=/app/config.yaml

# Run the Go app
ENTRYPOINT ["/odyscan"]
