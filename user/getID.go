package user

import (
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	"f-discover/services"

	"github.com/kataras/iris/v12"
)

type UserResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	CoverUrl     string `json:"coverUrl"`
	Job          string `json:"job"`
	AvatarUrl    string `json:"avatarUrl"`
	Quote        string `json:"quote"`
	Following    int    `json:"following"`
	Followers    int    `json:"followers"`
	FollowStatus int    `json:"followStatus"`
}

func GetID(ctx iris.Context) {
	usersCollection := services.GetInstance().StoreClient.Collection("users")

	userID := ctx.Params().Get("id")

	dsnap, err := usersCollection.Doc(userID).Get(instance.CtxBackground)
	if err != nil {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "User is inexistent"})
		return
	}
	var user models.User
	dsnap.DataTo(&user)

	var followStatus int

	currentUserID := helpers.GetCurrentUser(ctx).ID
	if currentUserID == "-1" {
		followStatus = -1
	} else if user.Followers[currentUserID] {
		followStatus = 1
	} else {
		followStatus = 0
	}

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: UserResponse{
			ID:           user.ID,
			Name:         user.Name,
			CoverUrl:     user.CoverUrl,
			AvatarUrl:    user.AvatarUrl,
			Job:          user.Job,
			Quote:        user.Quote,
			Following:    len(user.Following),
			Followers:    len(user.Followers),
			FollowStatus: followStatus,
		},
	})
}
