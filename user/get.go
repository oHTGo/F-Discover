package user

import (
	"context"
	"f-discover/firebase"
	"f-discover/interfaces"
	"f-discover/models"

	"github.com/kataras/iris/v12"
)

func Get(ctx iris.Context) {
	usersCollection := firebase.GetInstance().StoreClient.Collection("users")

	id := ctx.Values().GetString("id")

	dsnap, _ := usersCollection.Doc(id).Get(context.Background())
	var user models.User
	dsnap.DataTo(&user)

	res := Response{
		ID:        user.ID,
		Name:      user.Name,
		AvatarUrl: user.AvatarUrl,
		Following: len(user.Following),
		Followers: len(user.Followers),
	}

	ctx.JSON(interfaces.ISuccess{
		Message: "Get profile successfully",
		Data:    res,
	})
}
