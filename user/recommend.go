package user

import (
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/services"
	"math/rand"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kataras/iris/v12"
)

type RecommendUserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RecommendQuery struct {
	Max int `url:"max" json:"max"`
}

func Recommend(ctx iris.Context) {
	var query RecommendQuery
	if err := ctx.ReadQuery(&query); err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: "Query is invalid",
		})
		return
	}

	if errValidation := validation.ValidateStruct(&query,
		validation.Field(&query.Max, validation.Required),
	); errValidation != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFailWithErrors{Message: "Have validation error", Errors: errValidation})
		return
	}

	usersCollection := services.GetInstance().StoreClient.Collection("users")
	users, _ := usersCollection.Documents(instance.CtxBackground).GetAll()

	var recommendUsers []RecommendUserResponse

	if len(users) <= query.Max {
		for _, user := range users {
			id := user.Ref.ID
			name, _ := user.DataAt("name")
			recommendUsers = append(recommendUsers, RecommendUserResponse{
				ID:   id,
				Name: name.(string),
			})
		}
	} else {
		for _, position := range rand.Perm(len(users) - 1)[:query.Max] {
			id := users[position].Ref.ID
			name, _ := users[position].DataAt("name")
			recommendUsers = append(recommendUsers, RecommendUserResponse{
				ID:   id,
				Name: name.(string),
			})
		}
	}

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data:    recommendUsers,
	})
}
