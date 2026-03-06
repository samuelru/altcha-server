# Stage 1: Build the Go application
FROM golang:1.26.1-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum for dependency resolution
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# Use CGO_ENABLED=0 for a statically linked binary that runs on alpine/scratch
RUN CGO_ENABLED=0 GOOS=linux go build -v -o altcha-server .

# Stage 2: Final image
FROM alpine:3.19

# Install ca-certificates in case the server makes outgoing HTTPS requests
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/altcha-server .

# Ensure the application can write to the log/data files
# The app creates these files in the current working directory
RUN touch challenges.log solutions.txt verifications.log && \
    chmod 666 challenges.log solutions.txt verifications.log

# Expose the application port
EXPOSE 3947

# Define default environment variables
ENV ALTCHA_TTL="1h" \
    ALTCHA_CORS_ORIGIN="*" \
    IS_DEV="false"

# Run the application
CMD ["./altcha-server"]
