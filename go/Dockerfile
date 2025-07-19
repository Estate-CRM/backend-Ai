# Use a base Go image
FROM golang:1.24.4

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum first (optional but good for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go app
RUN go build -o server .

# Expose the port your app listens on
EXPOSE 8080

# Run the binary
CMD ["./server"]
