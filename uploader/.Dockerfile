FROM golang:1.23.2 AS builder

WORKDIR /app

COPY . .

# Build the Go application (binary output)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/uploader /app/cmd/uploader/main.go


######################### STAGE 2 #########################
# Use Node.js base image for the final container
FROM node:18

WORKDIR /app

COPY ./companion/package*.json ./

RUN npm install

COPY ./companion .

COPY .env ./

# Copy the binary from the builder image
COPY --from=builder /app/uploader /app/uploader

# Copy entrypoint script
COPY entry.sh /usr/local/bin/entry.sh
RUN chmod +x /usr/local/bin/entry.sh

EXPOSE 8081

# Run entrypoint script
CMD ["/usr/local/bin/entry.sh"]