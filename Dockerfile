# Use an official Golang runtime as a parent image
FROM golang:1.23.2-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy `go.mod` and `go.sum` files first for dependency caching
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Expose port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]