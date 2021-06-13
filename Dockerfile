FROM golang:1.15-alpine AS development

WORKDIR /app
COPY go.mod ./

RUN go mod download
COPY . .

RUN go mod tidy && go build -o app .

CMD go run main.go

FROM alpine:latest as production

WORKDIR /app
COPY . .
COPY --from=development /app/app .

CMD ./app