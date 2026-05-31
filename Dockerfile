FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./server.go
RUN CGO_ENABLED=0 GOOS=linux go build -o cli ./cmd/cli.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/cli .
EXPOSE 8000
CMD ["./server"]
