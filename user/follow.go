package user

import (
	"context"
	"f-discover/firebase"
	"f-discover/interfaces"
	"f-discover/models"
	"f-discover/params"

	"cloud.google.com/go/firestore"
	"github.com/kataras/iris/v12"
)

func Follow(ctx iris.Context) {
	usersCollection := firebase.GetInstance().StoreClient.Collection("users")

	var user params.UserID
	currentUserID := ctx.Values().GetString("id")

	if err := ctx.ReadParams(&user); err != nil {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "User ID is invalid"})
		return
	}

	if currentUserID == user.ID {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{Message: "Self-following is not allowed"})
		return
	}

	dsnapUser, err := usersCollection.Doc(user.ID).Get(context.Background())
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

	_, _ = usersCollection.Doc(user.ID).Update(context.Background(), []firestore.Update{
		{
			Path:  "followers." + currentUserID,
			Value: currentUserRef,
		},
	})

	_, _ = usersCollection.Doc(currentUserID).Update(context.Background(), []firestore.Update{
		{
			Path:  "following." + user.ID,
			Value: dsnapUser.Ref,
		},
	})

	ctx.JSON(interfaces.ISuccessNoData{
		Message: "Follow successfully",
	})
}
