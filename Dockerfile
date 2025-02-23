FROM golang:1.23.4 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy module files first for efficient caching
COPY odyscan/go.mod odyscan/go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY odyscan/ .

# Build the Go binary (targeting Linux)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o odyscan ./main.go

# ---- Run Stage (Minimal Image) ----
FROM gcr.io/distroless/static:latest

# Set the working directory inside the final image
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/odyscan .

# Copy any required static assets if needed
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates

# Define the entrypoint
ENTRYPOINT ["/root/odyscan"]
