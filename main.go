package main

import (
	"f-discover/authentication"
	"f-discover/env"
	"f-discover/errors"
	"f-discover/interfaces"
	"f-discover/location"
	"f-discover/logger"
	"f-discover/middlewares"
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

	app.HandleDir("/files", iris.Dir("./files"), iris.DirOptions{
		Cache: iris.DirCacheOptions{
			Enable: false,
		},
		Attachments: iris.Attachments{
			Enable: false,
		},
	})

	app.HandleDir("/location/images", iris.Dir("./location/images"), iris.DirOptions{
		Cache: iris.DirCacheOptions{
			Enable: false,
		},
		Attachments: iris.Attachments{
			Enable: false,
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

		userRouter.Get("/{id}", middlewares.CustomJWTMiddleware(j), user.GetID)

		userRouter.Get("/{id}/follow", j.Serve, user.CheckFollow)
		userRouter.Post("/{id}/follow", j.Serve, user.Follow)
		userRouter.Delete("/{id}/follow", j.Serve, user.Unfollow)

		userRouter.Get("/suggest", user.Suggest)
	}

	postRouter := api.Party("post")
	{
		postRouter.Post("/", j.Serve, post.Create)
		postRouter.Get("/{id}", middlewares.CustomJWTMiddleware(j), post.GetID)
		postRouter.Put("/{id}", j.Serve, post.Update)
		postRouter.Post("/{id}/upload-video", j.Serve, post.UploadVideo)
		postRouter.Delete("/{id}", j.Serve, post.Delete)

		postRouter.Get("/{id}/like", j.Serve, post.CheckLike)
		postRouter.Post("/{id}/like", j.Serve, post.Like)
		postRouter.Delete("/{id}/like", j.Serve, post.Unlike)

		postRouter.Post("/{id}/comment", j.Serve, post.CreatComment)
		postRouter.Get("/{id}/comment/", post.GetAllComment)
		postRouter.Get("/{id}/comment/{commentID}", post.GetComment)
		postRouter.Put("/{id}/comment/{commentID}", j.Serve, post.UpdateComment)
		postRouter.Delete("/{id}/comment/{commentID}", j.Serve, post.DeleteComment)

		postRouter.Get("/user/{id}", middlewares.CustomJWTMiddleware(j), post.GetListOfUser)
		postRouter.Get("/location/{id}", middlewares.CustomJWTMiddleware(j), post.GetListOfLocation)
		postRouter.Get("/following", j.Serve, post.GetListOfFollowing)
		postRouter.Get("/suggest", middlewares.CustomJWTMiddleware(j), post.Suggest)
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
