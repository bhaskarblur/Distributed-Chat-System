# Use Golang as the base image for building the application
FROM golang:1.23.1-bullseye AS builder

# Set the working directory
WORKDIR /build

# Set up Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o main ./cmd/server

# Use a smaller base image for the runtime
FROM debian:bullseye-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    libssl-dev zlib1g && \
    rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /build/main .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./main"]
