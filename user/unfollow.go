package user

import (
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	"f-discover/services"

	"cloud.google.com/go/firestore"
	"github.com/kataras/iris/v12"
)

func Unfollow(ctx iris.Context) {
	usersCollection := services.GetInstance().StoreClient.Collection("users")

	userID := ctx.Params().Get("id")
	currentUserID := ctx.Values().GetString("id")

	if currentUserID == userID {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{Message: "Self-following is not allowed"})
		return
	}

	dsnapUser, err := usersCollection.Doc(userID).Get(instance.CtxBackground)
	if err != nil {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "User is inexistent"})
		return
	}

	var userDoc models.User
	dsnapUser.DataTo(&userDoc)

	// Check current user has followed this user or not
	if userDoc.Followers[currentUserID] == nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{Message: "Current user has not followed this user"})
		return
	}

	_, _ = usersCollection.Doc(userID).Update(instance.CtxBackground, []firestore.Update{
		{
			Path:  "followers." + currentUserID,
			Value: firestore.Delete,
		},
	})

	_, _ = usersCollection.Doc(currentUserID).Update(instance.CtxBackground, []firestore.Update{
		{
			Path:  "following." + userID,
			Value: firestore.Delete,
		},
	})

	ctx.JSON(interfaces.ISuccessNoData{
		Message: "Unfollow successfully",
	})
}
