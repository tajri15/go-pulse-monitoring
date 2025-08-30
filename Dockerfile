# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-pulse-monitoring ./cmd/server

# Stage 2: Create a minimal final image
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /go-pulse-monitoring .

EXPOSE 8080

CMD ["./go-pulse-monitoring"]