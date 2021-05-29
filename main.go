package main

import (
	"f-discover/middlewares"
	"f-discover/user"
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

	userRouter := app.Party("user")
	{
		userRouter.Get("/", middlewares.SetAuthentication(), user.Get)
		userRouter.Put("/", middlewares.SetAuthentication(), user.UpdateProfile)
		userRouter.Post("/upload-avatar", middlewares.SetAuthentication(), user.UpdateAvatar)
		userRouter.Get("/{id}", user.GetID)
		userRouter.Post("/{id}/follow", middlewares.SetAuthentication(), user.Follow)
		userRouter.Post("/{id}/unfollow", middlewares.SetAuthentication(), user.Unfollow)
	}
	app.Listen(":" + envConfig["PORT"])
}

func helloWorldPoint(ctx iris.Context) {
	ctx.JSON("Hello World")
}
