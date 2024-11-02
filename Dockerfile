# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files for dependency resolution
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go binary for production
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Create a minimal final image
FROM alpine:latest

# Set environment variable to run Gin in production mode
ENV GIN_MODE=release


# Set the working directory inside the container
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/main .

COPY --from=builder /app/Reports ./Reports

# Expose the application port (optional, but not required for Cloud Run)
EXPOSE 8080

# Command to run the Go binary
CMD ["./main"]
