# ========================
# Build stage
# ========================
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o nanny-backend ./cmd/api


# ========================
# Runtime stage
# ========================
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/nanny-backend .

EXPOSE 8080

ENTRYPOINT ["./nanny-backend"]
