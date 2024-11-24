# Use official Golang image as a build stage
FROM golang:1.22.9-bullseye as builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the application binary targeting cmd/main.go
RUN go build -o distributed-chat-system ./cmd/main.go

# Use a smaller base image for the runtime stage
FROM alpine:latest

# Install certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/distributed-chat-system .

# Expose the application port
EXPOSE 8080

# Command to run the binary
CMD ["./distributed-chat-system"]
