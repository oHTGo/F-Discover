package post

import (
	"context"
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	IPost "f-discover/post/interfaces"
	"f-discover/services"
	"time"

	"cloud.google.com/go/firestore"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kataras/iris/v12"
)

type CreateCommentDTO struct {
	Content string `json:"content"`
}

type NewCommentResponse struct {
	ID        string       `json:"id"`
	Content   string       `json:"content"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	Author    IPost.Author `json:"author"`
}

func CreatComment(ctx iris.Context) {
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

	currentUser := helpers.GetCurrentUser(ctx)
	dsnap, _ := currentUser.Reference.Get(instance.CtxBackground)
	var author models.User
	dsnap.DataTo(&author)

	postID := ctx.Params().Get("id")
	_, err := postsCollection.Doc(postID).Get(instance.CtxBackground)
	if err != nil {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "Post is inexistent"})
		return
	}

	createdAt := time.Now()

	comment := models.Comment{
		ID:        helpers.RandomString(20),
		Content:   body.Content,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
		Author:    currentUser.Reference,
	}

	_, _ = postsCollection.Doc(postID).Update(context.Background(), []firestore.Update{
		{
			Path:  "comments." + comment.ID,
			Value: comment,
		},
	})

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: NewCommentResponse{
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