FROM golang:1.23.2 AS builder

WORKDIR /app

COPY ./manager/go.mod ./manager/go.sum /app/manager/

COPY ./go-common /app/go-common

WORKDIR /app/manager

RUN go mod download

COPY ./manager /app/manager

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main ./cmd

## Stage 2: Create a minimal runtime image
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
