package helpers

import (
	"f-discover/models"

	"github.com/kataras/iris/v12"
)

func GetFollowStatus(ctx iris.Context, user models.User) int {
	var followStatus int
	if GetCurrentUser(ctx).ID == "-1" || GetCurrentUser(ctx).ID == user.ID {
		followStatus = -1
	} else if user.Followers[GetCurrentUser(ctx).ID] {
		followStatus = 1
	} else {
		followStatus = 0
	}
	return followStatus
}
