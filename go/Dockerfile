# syntax=docker/dockerfile:1
FROM golang:1.24.4-alpine

# Install system dependencies
RUN apk add --no-cache git curl bash

# Install Air (live reload tool for Go)
RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air"]
