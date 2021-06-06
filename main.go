package main

import (
	"f-discover/authentication"
	"f-discover/env"
	"f-discover/post"
	"f-discover/user"

	"github.com/iris-contrib/middleware/cors"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

func main() {
	env.Get()

	app := iris.New()

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	app.UseRouter(crs)

	j := jwt.New(jwt.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(env.Get().JWT_SECRET), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	api := app.Party("/api/")

	api.Get("/", func(ctx iris.Context) { ctx.JSON("Hello World") })

	authenticationRouter := api.Party("authentication")
	{
		authenticationRouter.Post("/", authentication.ExchangeToken)
	}

	userRouter := api.Party("user", j.Serve)
	{
		userRouter.Get("/", user.Get)
		userRouter.Put("/", user.UpdateProfile)
		userRouter.Post("/upload-avatar", user.UploadAvatar)
		userRouter.Get("/{id}", user.GetID)
		userRouter.Post("/{id}/follow", user.Follow)
		userRouter.Post("/{id}/unfollow", user.Unfollow)
	}

	postRouter := api.Party("post", j.Serve)
	{
		postRouter.Post("/", post.Create)
		postRouter.Get("/{id}", post.GetID)
		postRouter.Post("/{id}/upload-files", post.UploadMediaFiles)

		postRouter.Get("/list/{id}", post.GetListOfUser)
	}

	app.Listen(":" + env.Get().PORT)
}
