# syntax=docker/dockerfile:1
FROM golang:1.24-alpine

RUN apk add --no-cache git curl

WORKDIR /app

# Install Air
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air"]
