package user

import (
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/services"

	"cloud.google.com/go/firestore"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kataras/iris/v12"
)

type UpdateProfileDTO struct {
	Name  string `json:"name"`
	Quote string `json:"quote"`
	Job   string `json:"job"`
}

type NewProfile struct {
	Name  string `json:"name"`
	Quote string `json:"quote"`
	Job   string `json:"job"`
}

func UpdateProfile(ctx iris.Context) {
	usersCollection := services.GetInstance().StoreClient.Collection("users")

	currenUser := helpers.GetCurrentUser(ctx)

	var body UpdateProfileDTO
	if err := ctx.ReadBody(&body); err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{Message: "Body is bad request"})
		return
	}

	if errValidation := validation.ValidateStruct(&body,
		validation.Field(&body.Name, validation.Required),
		validation.Field(&body.Quote, validation.Required),
		validation.Field(&body.Job, validation.Required),
	); errValidation != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFailWithErrors{Message: "Have validation error", Errors: errValidation})
		return
	}

	_, _ = usersCollection.Doc(currenUser.ID).Update(instance.CtxBackground, []firestore.Update{
		{
			Path:  "name",
			Value: string(body.Name),
		},
		{
			Path:  "quote",
			Value: string(body.Quote),
		},
		{
			Path:  "job",
			Value: string(body.Job),
		},
	})

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: NewProfile{
			Name:  string(body.Name),
			Quote: string(body.Quote),
			Job:   string(body.Job),
		},
	})
}
