package user

import (
	"context"
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	"f-discover/services"

	"cloud.google.com/go/firestore"
	"github.com/kataras/iris/v12"
)

func Follow(ctx iris.Context) {
	usersCollection := services.GetInstance().StoreClient.Collection("users")

	userID := ctx.Params().Get("id")
	currentUserID := helpers.GetCurrentUserID(ctx)

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
	if userDoc.Followers[currentUserID] != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{Message: "Current user has followed this user"})
		return
	}

	currentUserRef := usersCollection.Doc(currentUserID)

	_, _ = usersCollection.Doc(userID).Update(context.Background(), []firestore.Update{
		{
			Path:  "followers." + currentUserID,
			Value: currentUserRef,
		},
	})

	_, _ = usersCollection.Doc(currentUserID).Update(context.Background(), []firestore.Update{
		{
			Path:  "following." + userID,
			Value: dsnapUser.Ref,
		},
	})

	ctx.JSON(interfaces.ISuccessNoData{
		Message: "Success",
	})
}
