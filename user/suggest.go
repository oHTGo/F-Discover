package user

import (
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	"f-discover/services"
	"math/rand"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kataras/iris/v12"
)

type SuggestResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	CoverUrl     string `json:"coverUrl"`
	AvatarUrl    string `json:"avatarUrl"`
	Job          string `json:"job"`
	Quote        string `json:"quote"`
	FollowStatus int    `json:"followStatus"`
	Following    int    `json:"following"`
	Followers    int    `json:"followers"`
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
			var userData models.User
			user.DataTo(&userData)

			suggestedUsers = append(suggestedUsers, SuggestResponse{
				ID:           userData.ID,
				Name:         userData.Name,
				CoverUrl:     userData.CoverUrl,
				AvatarUrl:    userData.AvatarUrl,
				Job:          userData.Job,
				Quote:        userData.Quote,
				FollowStatus: helpers.GetFollowStatus(ctx, userData),
				Following:    len(userData.Following),
				Followers:    len(userData.Followers),
			})
		}
	} else {
		for _, position := range rand.Perm(len(users) - 1)[:query.Limit] {
			var userData models.User
			users[position].DataTo(&userData)

			suggestedUsers = append(suggestedUsers, SuggestResponse{
				ID:           userData.ID,
				Name:         userData.Name,
				CoverUrl:     userData.CoverUrl,
				AvatarUrl:    userData.AvatarUrl,
				Job:          userData.Job,
				Quote:        userData.Quote,
				FollowStatus: helpers.GetFollowStatus(ctx, userData),
				Following:    len(userData.Following),
				Followers:    len(userData.Followers),
			})
		}
	}

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data:    suggestedUsers,
	})
}
