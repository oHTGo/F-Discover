package post

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

func Like(ctx iris.Context) {
	postsCollection := services.GetInstance().StoreClient.Collection("posts")

	currentUser := helpers.GetCurrentUser(ctx)

	postID := ctx.Params().Get("id")
	dsnap, err := postsCollection.Doc(postID).Get(instance.CtxBackground)
	if err != nil {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "Post is inexistent"})
		return
	}

	var post models.Post
	dsnap.DataTo(&post)

	// Check current user has liked this post or not
	if post.Likes[currentUser.ID] {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{Message: "Current user has liked this post"})
		return
	}

	_, _ = postsCollection.Doc(postID).Update(context.Background(), []firestore.Update{
		{
			Path:  "likes." + currentUser.ID,
			Value: true,
		},
	})

	ctx.JSON(interfaces.ISuccessNoData{
		Message: "Success",
	})
}
