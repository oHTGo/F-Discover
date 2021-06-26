package user

import (
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/services"
	"math/rand"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kataras/iris/v12"
)

type SuggestResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatarUrl"`
}

type SuggestQuery struct {
	Limit int `url:"limit" json:"limit"`
}

func Suggest(ctx iris.Context) {
	var query SuggestQuery
	if err := ctx.ReadQuery(&query); err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: "Query is invalid",
		})
		return
	}

	if errValidation := validation.ValidateStruct(&query,
		validation.Field(&query.Limit, validation.Required, validation.Min(1)),
	); errValidation != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFailWithErrors{Message: "Have validation error", Errors: errValidation})
		return
	}

	usersCollection := services.GetInstance().StoreClient.Collection("users")
	users, _ := usersCollection.Documents(instance.CtxBackground).GetAll()

	var suggestedUsers []SuggestResponse

	if len(users) <= query.Limit {
		for _, user := range users {
			id := user.Ref.ID
			name, _ := user.DataAt("name")
			avatarUrl, _ := user.DataAt("avatarUrl")
			suggestedUsers = append(suggestedUsers, SuggestResponse{
				ID:        id,
				Name:      name.(string),
				AvatarUrl: avatarUrl.(string),
			})
		}
	} else {
		for _, position := range rand.Perm(len(users) - 1)[:query.Limit] {
			id := users[position].Ref.ID
			name, _ := users[position].DataAt("name")
			avatarUrl, _ := users[position].DataAt("avatarUrl")
			suggestedUsers = append(suggestedUsers, SuggestResponse{
				ID:        id,
				Name:      name.(string),
				AvatarUrl: avatarUrl.(string),
			})
		}
	}

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data:    suggestedUsers,
	})
}
