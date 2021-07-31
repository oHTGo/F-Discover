package helpers

import (
	"f-discover/models"

	"github.com/kataras/iris/v12"
)

func GetLikeStatus(ctx iris.Context, post models.Post) int {
	var likeStatus int
	if GetCurrentUser(ctx).ID == "-1" {
		likeStatus = -1
	} else if post.Likes[GetCurrentUser(ctx).ID] {
		likeStatus = 1
	} else {
		likeStatus = 0
	}

	return likeStatus
}
