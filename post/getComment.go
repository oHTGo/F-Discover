package post

import (
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	IPost "f-discover/post/interfaces"
	"f-discover/services"

	"github.com/kataras/iris/v12"
)

func GetComment(ctx iris.Context) {
	postsCollection := services.GetInstance().StoreClient.Collection("posts")

	postID := ctx.Params().Get("id")
	dsnap, err := postsCollection.Doc(postID).Get(instance.CtxBackground)
	if err != nil {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "Post is inexistent"})
		return
	}
	var post models.Post
	dsnap.DataTo(&post)

	commentID := ctx.Params().Get("commentID")
	if _, ok := post.Comments[commentID]; !ok {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "Comment is inexistent"})
		return
	}

	var comment models.Comment = post.Comments[commentID]

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: IPost.CommentWithoutAuthor{
			ID:        comment.ID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
		},
	})
}
