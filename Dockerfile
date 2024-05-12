FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN apt-get update && apt-get install -y libvips-dev

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/main ./cmd/

ENV GIN_MODE=release
ENV ENV="dev"
ENV SHUTDOWN_TIMEOUT="10s"
ENV JWT_SECRET="some_hard_secret"

ENV HTTP_ADDRESS="localhost:5000"
ENV HTTP_TIMEOUT="4s"
ENV HTTP_IDLE_TIMEOUT="4s"

ENV USER_SERVICE_ADDRESS="localhost:44044"
ENV USER_SERVICE_TIMEOUT="3s"
ENV USER_SERVICE_RETRIES_COUNT=3
ENV CLUB_SERVICE_ADDRESS="localhost:44045"
ENV CLUB_SERVICE_TIMEOUT="3s"
ENV CLUB_SERVICE_RETRIES_COUNT=3

ENV AWS_REGION="us-east-1"
ENV AWS_ACCESS_KEY_ID="access_key_id"
ENV AWS_SECRET_ACCESS_KEY="secret_access_key"

EXPOSE 5000

ENTRYPOINT ["./build/main"]