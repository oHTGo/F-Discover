package helpers

import (
	"f-discover/interfaces"
	"f-discover/logger"
	"f-discover/services"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
)

func GetCurrentUser(ctx iris.Context) *interfaces.CurrentUser {
	var id string

	if ok := ctx.Values().Get("jwt"); ok != nil {
		id = ok.(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)
		logger.Info("Access log", "ID user: "+id)
	} else {
		id = "-1"
	}

	return &interfaces.CurrentUser{
		ID:        id,
		Reference: services.GetInstance().StoreClient.Collection("users").Doc(id),
	}
}
