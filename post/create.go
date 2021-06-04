package post

import (
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	IPost "f-discover/post/interfaces"
	"f-discover/services"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kataras/iris/v12"
)

type CreatePostDTO struct {
	Content string `json:"content"`
}

type NewPost struct {
	ID      string       `json:"id"`
	Content string       `json:"content"`
	Author  IPost.Author `json:"author"`
}

func Create(ctx iris.Context) {
	postsCollection := services.GetInstance().StoreClient.Collection("posts")
	usersCollection := services.GetInstance().StoreClient.Collection("users")

	userID := helpers.GetCurrentUserID(ctx)
	authorRef := usersCollection.Doc(userID)
	var author models.User
	dsnap, _ := authorRef.Get(instance.CtxBackground)
	dsnap.DataTo(&author)

	var body CreatePostDTO
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

	newPost := postsCollection.NewDoc()
	post := models.Post{
		ID:      newPost.ID,
		Content: body.Content,
		Author:  authorRef,
	}
	_, _ = newPost.Set(instance.CtxBackground, post)

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: NewPost{
			ID:      post.ID,
			Content: post.Content,
			Author: IPost.Author{
				ID:        author.ID,
				Name:      author.Name,
				AvatarUrl: author.AvatarUrl,
			},
		},
	})
}
