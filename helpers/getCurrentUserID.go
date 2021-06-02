package helpers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
)

func GetCurrentUserID(ctx iris.Context) string {
	user := ctx.Values().Get("jwt").(*jwt.Token)
	return user.Claims.(jwt.MapClaims)["id"].(string)
}
