FROM golang:1.15-buster AS development

WORKDIR /app
COPY go.mod ./

RUN go mod download
COPY . .

RUN go mod tidy

# Builder for production
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o app .

CMD go run main.go

FROM alpine:latest as production

WORKDIR /app
COPY . .
COPY --from=development /app/app .

CMD ./app