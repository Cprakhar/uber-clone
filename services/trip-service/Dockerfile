# Build stage
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the trip-service binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./services/trip-service

# Final stage
FROM alpine:latest

# Install ca-certificates
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port 9000
EXPOSE 9000

# Run the binary
CMD ["./main"]