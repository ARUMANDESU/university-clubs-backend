FROM golang:1.21 as builder

WORKDIR /app

# Copy the go.mod and go.sum files first and download the dependencies.
# This is done separately from copying the entire source code to leverage Docker cache
# and avoid re-downloading dependencies if they haven't changed.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application's source code.
COPY . .

# Build the application. This assumes you have a main package at the root of your project.
# Adjust the path to the main package if it's located elsewhere.
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/main ./cmd/

ENV GIN_MODE=release
# Expose the port your application listens on.
EXPOSE 5000

CMD ["./build/main", "--config=./config/dev.yaml"]