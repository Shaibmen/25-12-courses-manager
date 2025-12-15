FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd
RUN go build -o /app/main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/internal/documents ./internal/documents
COPY .env .

CMD ["./main"]
