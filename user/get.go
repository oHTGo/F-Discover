package user

import (
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"

	"github.com/kataras/iris/v12"
)

type CurrentUser struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CoverUrl  string `json:"coverUrl"`
	Job       string `json:"job"`
	AvatarUrl string `json:"avatarUrl"`
	Quote     string `json:"quote"`
	Following int    `json:"following"`
	Followers int    `json:"followers"`
}

func Get(ctx iris.Context) {
	currentUser := helpers.GetCurrentUser(ctx)

	dsnap, err := currentUser.Reference.Get(instance.CtxBackground)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, interfaces.IFail{Message: "Get profile failed"})
		return
	}

	var user models.User
	dsnap.DataTo(&user)

	res := CurrentUser{
		ID:        user.ID,
		Name:      user.Name,
		CoverUrl:  user.CoverUrl,
		AvatarUrl: user.AvatarUrl,
		Job:       user.Job,
		Quote:     user.Quote,
		Following: len(user.Following),
		Followers: len(user.Followers),
	}

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data:    res,
	})
}
