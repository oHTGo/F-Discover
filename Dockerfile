FROM golang:1.15-buster AS development

RUN apt update
RUN apt install -y ffmpeg

WORKDIR /app
COPY go.mod ./

RUN go mod download
COPY . .

RUN go mod tidy

# Builder for production
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o server .

CMD go run main.go

FROM alpine:latest as production

RUN apk add ffmpeg

WORKDIR /app
COPY . .
COPY --from=development /app/server .

CMD ./server