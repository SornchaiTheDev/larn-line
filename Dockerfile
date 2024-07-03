# First stage: Build the Go application
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o larn ./cmd

# Second stage: Create a small image to run the Go application
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/larn .

# Expose the application port
EXPOSE 3000

# Run the binary
CMD ["./larn"]
