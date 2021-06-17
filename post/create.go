package post

import (
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/location"
	"f-discover/models"
	IPost "f-discover/post/interfaces"
	"f-discover/services"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kataras/iris/v12"
)

type CreatePostDTO struct {
	Content  string `json:"content"`
	Location string `json:"location"`
}

type NewPost struct {
	ID        string       `json:"id"`
	Content   string       `json:"content"`
	Author    IPost.Author `json:"author"`
	Location  string       `json:"location"`
	CreatedAt time.Time    `json:"createdAt"`
}

func Create(ctx iris.Context) {
	postsCollection := services.GetInstance().StoreClient.Collection("posts")

	var body CreatePostDTO
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

	currentUser := helpers.GetCurrentUser(ctx)

	dsnap, _ := currentUser.Reference.Get(instance.CtxBackground)
	var user models.User
	dsnap.DataTo(&user)

	createdAt := time.Now()

	newPost := postsCollection.NewDoc()
	post := models.Post{
		ID:        newPost.ID,
		Content:   body.Content,
		Author:    currentUser.Reference,
		Location:  body.Location,
		CreatedAt: createdAt,
	}
	_, _ = newPost.Set(instance.CtxBackground, post)

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: NewPost{
			ID:        post.ID,
			Content:   post.Content,
			Location:  location.GetName(post.Location),
			CreatedAt: post.CreatedAt,
			Author: IPost.Author{
				ID:        user.ID,
				Name:      user.Name,
				AvatarUrl: user.AvatarUrl,
				Job:       user.Job,
			},
		},
	})
}
