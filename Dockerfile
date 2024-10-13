# Stage 1: Build the Go application
FROM golang:1.23.2-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go app
RUN go build -o jwe-go .

# Stage 2: Create a smaller image to run the Go application
FROM alpine:latest

# Install certificates (if your app needs HTTPS calls)
# RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/jwe-go .

# Expose the port the app runs on (if needed)
EXPOSE 8080

# Command to run the executable
CMD ["./jwe-go"]
