# Build Stage
FROM golang:1.23.4 AS builder

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum first (for dependency caching)
COPY odyscan/go.mod ./
COPY odyscan/go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o odyscan .odyscan/cmd/main.go

# Final Image
FROM gcr.io/distroless/static:latest

# Set working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/odyscan .

# Set entrypoint
ENTRYPOINT ["/root/odyscan"]
