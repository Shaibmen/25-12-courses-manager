FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd
RUN go build -o /app/main .

FROM alpine:latest

RUN apk add --no-cache python3 py3-pip

COPY scripts/ /app/scripts/
COPY internal/count.json /app/

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/internal/documents ./internal/documents

RUN pip3 install --break-system-packages --no-cache-dir python-docx

CMD ["./main"]
