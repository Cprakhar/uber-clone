# Build stage
FROM golang:1.25.1-alpine AS builder
# Set working directory
WORKDIR /app
# Copy go mod and sum files
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download
# Copy the source code
COPY services/api-gateway/ ./services/api-gateway/
COPY shared/ ./shared/
# Build the api-gateway binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./services/api-gateway

# Final stage
FROM alpine:latest
# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates
# Set working directory
WORKDIR /root/
# Copy the binary from builder stage
COPY --from=builder /app/main .
# Expose port 8080
EXPOSE 8080
# Run the binary
CMD ["./main"]