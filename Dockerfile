FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o previewer ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/previewer /app/previewer

EXPOSE 8080

ENTRYPOINT ["/app/previewer"]
