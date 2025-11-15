FROM golang:1.22-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

COPY go.mod go.sum* ./
RUN go mod download || true

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /build/api .
COPY --from=builder /build/internal/seed ./internal/seed

EXPOSE 8083

CMD ["./api"]
