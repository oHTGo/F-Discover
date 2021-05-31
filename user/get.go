package user

import (
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	"f-discover/services"
	IUser "f-discover/user/interfaces"

	"github.com/kataras/iris/v12"
)

func Get(ctx iris.Context) {
	usersCollection := services.GetInstance().StoreClient.Collection("users")

	id := ctx.Values().GetString("id")

	dsnap, err := usersCollection.Doc(id).Get(instance.CtxBackground)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, interfaces.IFail{Message: "Get profile failed"})
		return
	}

	var user models.User
	dsnap.DataTo(&user)

	res := IUser.Info{
		ID:        user.ID,
		Name:      user.Name,
		AvatarUrl: user.AvatarUrl,
		Following: len(user.Following),
		Followers: len(user.Followers),
	}

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data:    res,
	})
}
