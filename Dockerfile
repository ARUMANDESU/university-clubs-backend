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

ENV GIN_MODE=release\
    ENV="dev"\
    SHUTDOWN_TIMEOUT="10s"\
    HTTP_ADDRESS="localhost:5000"\
    HTTP_TIMEOUT="4s"\
    HTTP_IDLE_TIMEOUT="4s"\
    USER_SERVICE_ADDRESS="localhost:44044"\
    USER_SERVICE_TIMEOUT="3s"\
    USER_SERVICE_RETRIES_COUNT=3\
    CLUB_SERVICE_ADDRESS="localhost:44045"\
    CLUB_SERVICE_TIMEOUT="3s"\
    CLUB_SERVICE_RETRIES_COUNT=3

# Expose the port your application listens on.
EXPOSE 5000

CMD ["./build/main"]