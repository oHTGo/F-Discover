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

func DeleteComment(ctx iris.Context) {
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

	commentID := ctx.Params().Get("commentID")
	if _, ok := post.Comment[commentID]; !ok {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "Comment is inexistent"})
		return
	}

	var comment models.Comment = post.Comment[commentID]

	if comment.Author.ID != currentUser.ID {
		ctx.StopWithJSON(iris.StatusForbidden, interfaces.IFail{Message: "User is not the author of the comment"})
		return
	}

	_, _ = postsCollection.Doc(postID).Update(context.Background(), []firestore.Update{
		{
			Path:  "comments." + comment.ID,
			Value: firestore.Delete,
		},
	})

	ctx.JSON(interfaces.ISuccessNoData{
		Message: "Success",
	})
}
