package middlewares

import (
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

func CustomJWTMiddleware(j *jwt.Middleware) iris.Handler {
	return func(ctx iris.Context) {
		j.CheckJWT(ctx)
		ctx.Next()
	}
}
