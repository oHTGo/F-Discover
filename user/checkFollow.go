package user

import (
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	"f-discover/services"

	"github.com/kataras/iris/v12"
)

type CheckFollowResponse struct {
	Followed bool `json:"followed"`
}

func CheckFollow(ctx iris.Context) {
	usersCollection := services.GetInstance().StoreClient.Collection("users")

	userID := ctx.Params().Get("id")
	currentUser := helpers.GetCurrentUser(ctx)

	if currentUser.ID == userID {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{Message: "Self-following is not allowed"})
		return
	}

	dsnapUser, err := usersCollection.Doc(userID).Get(instance.CtxBackground)
	if err != nil {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "User is inexistent"})
		return
	}

	var user models.User
	dsnapUser.DataTo(&user)

	var followed bool = false

	// Check current user has followed this user or not
	if _, ok := user.Followers[currentUser.ID]; ok {
		followed = true
	}

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: CheckFollowResponse{
			Followed: followed,
		},
	})
}
