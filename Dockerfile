FROM golang:1.16-buster

WORKDIR /app
COPY . ./

RUN go mod download && \
    go mod tidy

CMD go run main.go