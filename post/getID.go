package post

import (
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/location"
	"f-discover/models"
	IPost "f-discover/post/interfaces"
	"f-discover/services"

	"github.com/kataras/iris/v12"
)

func GetID(ctx iris.Context) {
	postsCollection := services.GetInstance().StoreClient.Collection("posts")
	usersCollection := services.GetInstance().StoreClient.Collection("users")

	postID := ctx.Params().Get("id")
	dsnap, err := postsCollection.Doc(postID).Get(instance.CtxBackground)
	if err != nil {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "Post is inexistent"})
		return
	}

	var post models.Post
	dsnap.DataTo(&post)

	dsnap, _ = usersCollection.Doc(post.Author.ID).Get(instance.CtxBackground)
	var author models.User
	dsnap.DataTo(&author)

	var likeStatus int
	if helpers.GetCurrentUser(ctx).ID == "-1" {
		likeStatus = -1
	} else if post.Likes[helpers.GetCurrentUser(ctx).ID] {
		likeStatus = 1
	} else {
		likeStatus = 0
	}

	res := IPost.Info{
		ID:           post.ID,
		Content:      post.Content,
		ThumbnailUrl: post.ThumbnailUrl,
		VideoUrl:     post.VideoUrl,
		Likes:        len(post.Likes),
		LikeStatus:   likeStatus,
		Comments:     len(post.Comments),
		Location:     location.GetName(post.Location),
		CreatedAt:    post.CreatedAt,
		Author: IPost.Author{
			ID:        author.ID,
			Name:      author.Name,
			AvatarUrl: author.AvatarUrl,
			Job:       author.Job,
		},
	}

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data:    res,
	})
}
