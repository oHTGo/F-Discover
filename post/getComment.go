package post

import (
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	IPost "f-discover/post/interfaces"
	"f-discover/services"
	"time"

	"github.com/kataras/iris/v12"
)

type GetCommentResponse struct {
	ID        string       `json:"id"`
	Content   string       `json:"content"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	Author    IPost.Author `json:"author"`
}

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

	dsnap, _ = comment.Author.Get(instance.CtxBackground)
	var author models.User
	dsnap.DataTo(&author)

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: GetCommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			Author: IPost.Author{
				ID:        author.ID,
				Name:      author.Name,
				AvatarUrl: author.AvatarUrl,
			},
		},
	})
}
