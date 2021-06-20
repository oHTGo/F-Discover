package post

import (
	"context"
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	IPost "f-discover/post/interfaces"
	"f-discover/services"

	"cloud.google.com/go/firestore"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kataras/iris/v12"
)

type UpdateCommentDTO struct {
	Content string `json:"content"`
}

func UpdateComment(ctx iris.Context) {
	postsCollection := services.GetInstance().StoreClient.Collection("posts")

	var body CreateCommentDTO
	if err := ctx.ReadBody(&body); err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{Message: "Body is bad request"})
		return
	}

	if errValidation := validation.ValidateStruct(&body,
		validation.Field(&body.Content, validation.Required),
	); errValidation != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFailWithErrors{Message: "Have validation error", Errors: errValidation})
		return
	}

	helpers.EscapeString(&body)

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
	if _, ok := post.Comments[commentID]; !ok {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "Comment is inexistent"})
		return
	}

	var comment models.Comment = post.Comments[commentID]

	if comment.Author.ID != currentUser.ID {
		ctx.StopWithJSON(iris.StatusForbidden, interfaces.IFail{Message: "User is not the author of the comment"})
		return
	}

	comment.Content = body.Content

	_, _ = postsCollection.Doc(postID).Update(context.Background(), []firestore.Update{
		{
			Path:  "comments." + comment.ID,
			Value: comment,
		},
	})

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: IPost.CommentWithoutAuthor{
			ID:        comment.ID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
		},
	})
}
