package post

import (
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	"f-discover/services"

	"github.com/kataras/iris/v12"
)

type CheckLikeResponse struct {
	Liked bool `json:"liked"`
}

func CheckLike(ctx iris.Context) {
	postsCollection := services.GetInstance().StoreClient.Collection("posts")

	currentUser := helpers.GetCurrentUser(ctx)

	postID := ctx.Params().Get("id")
	dsnap, err := postsCollection.Doc(postID).Get(instance.CtxBackground)
	if err != nil {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "Post is inexistent"})
		return
	}

	var post models.Post
	dsnap.DataTo(&post)

	var liked bool = false
	// Check current user has liked this post or not
	if post.Likes[currentUser.ID] {
		liked = true
	}

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: CheckLikeResponse{
			Liked: liked,
		},
	})
}
