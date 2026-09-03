# ---- Build stage ----

FROM golang:1.26.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .
RUN swag init -g cmd/api/main.go -o docs
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api/main.go

# ---- Run stage ----

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./server"]
