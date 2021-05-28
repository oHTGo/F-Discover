package user

import (
	"context"
	"f-discover/firebase"
	"f-discover/interfaces"
	"f-discover/models"

	"github.com/kataras/iris/v12"
)

func GetID(ctx iris.Context) {
	usersCollection := firebase.GetInstance().StoreClient.Collection("users")

	userID := ctx.Params().Get("id")
	// if err := ctx.ReadParams(&param); err != nil {
	// 	ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "User ID is invalid"})
	// 	return
	// }

	dsnap, err := usersCollection.Doc(userID).Get(context.Background())
	if err != nil {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "User is inexistent"})
		return
	}
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
		Message: "Get user successfully",
		Data:    res,
	})
}
