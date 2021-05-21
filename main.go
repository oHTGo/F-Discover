package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
)

func main() {
	envConfig, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := iris.New()
	app.Get("/", helloWorldPoint)

	app.Listen(":" + envConfig["PORT"])
}

func helloWorldPoint(ctx iris.Context) {
	ctx.JSON("Hello World")
}
