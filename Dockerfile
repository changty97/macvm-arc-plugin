# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o macvm-plugin .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/macvm-plugin .
ENTRYPOINT ["./macvm-plugin"]
