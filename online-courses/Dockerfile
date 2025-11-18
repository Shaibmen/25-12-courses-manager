FROM golang:1.24.4-alpine as builder
RUN apk add --no-cache postgresql-client
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/app
RUN go build -o /app/main .


FROM alpine:latest
RUN apk add --no-cache postgresql-client
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]