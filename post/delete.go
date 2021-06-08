package post

import (
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	"f-discover/services"

	"github.com/kataras/iris/v12"
)

func Delete(ctx iris.Context) {
	postsCollection := services.GetInstance().StoreClient.Collection("posts")
	usersCollection := services.GetInstance().StoreClient.Collection("users")

	postID := ctx.Params().Get("id")
	dsnap, err := postsCollection.Doc(postID).Get(instance.CtxBackground)
	if err != nil {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "Post is inexistent"})
		return
	}
	var post models.Post
	dsnap.DataTo(&post)

	userID := helpers.GetCurrentUserID(ctx)
	authorRef := usersCollection.Doc(userID)

	if post.Author.ID != authorRef.ID {
		ctx.StopWithJSON(iris.StatusForbidden, interfaces.IFail{Message: "User is not the author of the post"})
		return
	}

	_, _ = postsCollection.Doc(postID).Delete(instance.CtxBackground)

	ctx.JSON(interfaces.ISuccessNoData{
		Message: "Success",
	})
}
