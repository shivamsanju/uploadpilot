# Stage 1: Build the Go application
FROM golang:1.23.2-alpine as builder

# Set the working directory in the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go dependencies
RUN go mod tidy

# Copy the rest of the application files
COPY ./go .

# Build the Go application (binary output)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server /app/cmd/main.go


######################### STAGE 2 #########################
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder image
COPY --from=builder /app/server .
COPY --from=builder /app/.env .

# Expose the port the app will run on
EXPOSE 8081

# Command to run the binary
CMD ["./server"]
