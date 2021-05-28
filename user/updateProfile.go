package user

import (
	"context"
	"f-discover/firebase"
	"f-discover/interfaces"

	"cloud.google.com/go/firestore"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kataras/iris/v12"
)

type ProfileBody struct {
	Name string `json:"name"`
}

func UpdateProfile(ctx iris.Context) {
	usersCollection := firebase.GetInstance().StoreClient.Collection("users")

	id := ctx.Values().GetString("id")

	var body ProfileBody
	if err := ctx.ReadBody(&body); err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{Message: "Body is bad request"})
		return
	}

	if errValidation := validation.ValidateStruct(&body,
		validation.Field(&body.Name, validation.Required),
	); errValidation != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFailWithErrors{Message: "Have validation error", Errors: errValidation})
		return
	}

	_, _ = usersCollection.Doc(id).Update(context.Background(), []firestore.Update{
		{
			Path:  "name",
			Value: string(body.Name),
		},
	})

	ctx.JSON(interfaces.ISuccessNoData{
		Message: "Update profile successfully",
	})
}
