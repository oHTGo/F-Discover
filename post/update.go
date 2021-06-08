package post

import (
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	"f-discover/services"

	"cloud.google.com/go/firestore"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kataras/iris/v12"
)

type UpdatePostDTO struct {
	Content string `json:"content"`
}

type NewUpdatePost struct {
	Content string `json:"content"`
}

func Update(ctx iris.Context) {
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
	var author models.User
	dsnap, _ = authorRef.Get(instance.CtxBackground)
	dsnap.DataTo(&author)

	if post.Author.ID != authorRef.ID {
		ctx.StopWithJSON(iris.StatusForbidden, interfaces.IFail{Message: "User is not the author of the post"})
		return
	}

	var body UpdatePostDTO
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

	_, _ = postsCollection.Doc(postID).Update(instance.CtxBackground, []firestore.Update{
		{
			Path:  "content",
			Value: body.Content,
		},
	})

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: NewUpdatePost{
			Content: body.Content,
		},
	})
}
