package main

import (
	"f-discover/authentication"
	"f-discover/env"
	"f-discover/errors"
	"f-discover/interfaces"
	"f-discover/location"
	"f-discover/logger"
	"f-discover/post"
	"f-discover/search"
	"f-discover/user"
	"strings"

	"github.com/iris-contrib/middleware/cors"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

func main() {
	env.Get()

	logger.Init()

	app := iris.New()
	app.Use(errors.Handle())
	app.Use(logger.Handle())

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
		ErrorHandler: func(ctx iris.Context, err error) {
			ctx.StopExecution()
			ctx.StatusCode(iris.StatusUnauthorized)

			errMessage := []rune(err.Error())
			message := strings.ToUpper(string(errMessage[:1])) + strings.ToLower(string(errMessage[1:]))
			ctx.JSON(interfaces.IFail{
				Message: string(message),
			})
		},
	})

	api := app.Party("/api/")

	api.Get("/", func(ctx iris.Context) { ctx.JSON("Hello World") })

	authenticationRouter := api.Party("authentication")
	{
		authenticationRouter.Post("/", authentication.ExchangeToken)
	}

	userRouter := api.Party("user")
	{
		userRouter.Get("/", j.Serve, user.Get)
		userRouter.Put("/", j.Serve, user.UpdateProfile)
		userRouter.Post("/upload-avatar", j.Serve, user.UploadAvatar)
		userRouter.Post("/upload-cover", j.Serve, user.UploadCover)
		userRouter.Get("/{id}", user.GetID)
		userRouter.Post("/{id}/follow", j.Serve, user.Follow)
		userRouter.Post("/{id}/unfollow", j.Serve, user.Unfollow)

		userRouter.Get("/suggest", user.Suggest)
	}

	postRouter := api.Party("post", j.Serve)
	{
		postRouter.Post("/", post.Create)
		postRouter.Get("/{id}", post.GetID)
		postRouter.Put("/{id}", post.Update)
		postRouter.Post("/{id}/upload-video", post.UploadVideo)
		postRouter.Post("/{id}/like", post.Like)
		postRouter.Post("/{id}/unlike", post.Unlike)
		postRouter.Delete("/{id}", post.Delete)

		postRouter.Post("/{id}/comment", post.CreatComment)
		postRouter.Get("/{id}/comment/{commentID}", post.GetComment)
		postRouter.Put("/{id}/comment/{commentID}", post.UpdateComment)
		postRouter.Delete("/{id}/comment/{commentID}", post.DeleteComment)

		postRouter.Get("/user/{id}", post.GetListOfUser)
		postRouter.Get("/location/{id}", post.GetListOfLocation)
		postRouter.Get("/following", post.GetListOfFollowing)
		postRouter.Get("/suggest", post.Suggest)
	}

	locationRouter := api.Party("location")
	{
		locationRouter.Get("/", location.GetList)
	}

	searchRouter := api.Party("search")
	{
		searchRouter.Get("/", search.Search)
	}

	app.Listen(":" + env.Get().PORT)
}
