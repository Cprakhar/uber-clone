# Build Stage
FROM golang:1.25.1-bookworm AS builder
# Install dependencies
RUN apt-get update && apt-get install -y pkg-config build-essential librdkafka-dev
# Set working directory
WORKDIR /app
# Copy go mod and sum files
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download
# Copy the source code
COPY services/driver-service/ ./services/driver-service/
COPY shared/ ./shared/
# Build the driver-service binary
RUN CGO_ENABLED=1 GOOS=linux CC=gcc go build -a -installsuffix cgo -o main ./services/driver-service

# Final stage
FROM debian:bookworm-slim
# Install dependencies
RUN apt-get update && apt-get install -y librdkafka1 ca-certificates
# Set working directory
WORKDIR /root/
# Copy the binary from builder stage
COPY --from=builder /app/main .
# Expose port 9100
EXPOSE 9100
# Run the binary
CMD ["./main"]