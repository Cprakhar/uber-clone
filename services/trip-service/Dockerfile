# Build stage
FROM golang:1.25.1-alpine AS builder
# Set working directory
WORKDIR /app
# Copy go mod and sum files
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download
# Copy the source code
COPY services/trip-service/ ./services/trip-service/
COPY shared/ ./shared/
# Build the trip-service binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./services/trip-service

# Final stage
FROM alpine:latest
# Install ca-certificates
RUN apk --no-cache add ca-certificates
# Set working directory
WORKDIR /root/
# Copy the binary from builder stage
COPY --from=builder /app/main .
# Expose port 9000
EXPOSE 9000
# Run the binary
CMD ["./main"]