package post

import (
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/location"
	"f-discover/models"
	"f-discover/services"

	"cloud.google.com/go/firestore"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kataras/iris/v12"
)

type UpdatePostDTO struct {
	Content  string `json:"content"`
	Location string `json:"location"`
}

type NewUpdatePost struct {
	Content  string `json:"content"`
	Location string `json:"location"`
}

func Update(ctx iris.Context) {
	postsCollection := services.GetInstance().StoreClient.Collection("posts")

	var body UpdatePostDTO
	if err := ctx.ReadBody(&body); err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{Message: "Body is bad request"})
		return
	}

	if errValidation := validation.ValidateStruct(&body,
		validation.Field(&body.Content, validation.Required),
		validation.Field(&body.Location, validation.Required, validation.By(location.CheckValidationForValidator)),
	); errValidation != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFailWithErrors{Message: "Have validation error", Errors: errValidation})
		return
	}

	postID := ctx.Params().Get("id")
	dsnap, err := postsCollection.Doc(postID).Get(instance.CtxBackground)
	if err != nil {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "Post is inexistent"})
		return
	}
	var post models.Post
	dsnap.DataTo(&post)

	currentUser := helpers.GetCurrentUser(ctx)

	if post.Author.ID != currentUser.ID {
		ctx.StopWithJSON(iris.StatusForbidden, interfaces.IFail{Message: "User is not the author of the post"})
		return
	}

	_, _ = postsCollection.Doc(postID).Update(instance.CtxBackground, []firestore.Update{
		{
			Path:  "content",
			Value: body.Content,
		},
		{
			Path:  "location",
			Value: body.Location,
		},
	})

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: NewUpdatePost{
			Content:  body.Content,
			Location: location.GetName(body.Location),
		},
	})
}
