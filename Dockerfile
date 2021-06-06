FROM golang:1.15-alpine

WORKDIR /app
COPY . ./

RUN go mod download && \
    go mod tidy

CMD go run main.go