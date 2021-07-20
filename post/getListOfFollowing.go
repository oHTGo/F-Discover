package post

import (
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/location"
	"f-discover/models"
	PostHelpers "f-discover/post/helpers"
	IPost "f-discover/post/interfaces"
	"f-discover/services"

	"cloud.google.com/go/firestore"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kataras/iris/v12"
)

type GetListOfFollowingQuery struct {
	Page  int `url:"page" json:"page"`
	Limit int `url:"limit" json:"limit"`
}

func GetListOfFollowing(ctx iris.Context) {
	postsCollection := services.GetInstance().StoreClient.Collection("posts")
	usersCollection := services.GetInstance().StoreClient.Collection("users")

	var query GetListOfFollowingQuery
	if err := ctx.ReadQuery(&query); err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: "Query is invalid",
		})
		return
	}
	if errValidation := validation.ValidateStruct(&query,
		validation.Field(&query.Page, validation.Required, validation.Min(1)),
		validation.Field(&query.Limit, validation.Required, validation.Min(1)),
	); errValidation != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFailWithErrors{Message: "Have validation error", Errors: errValidation})
		return
	}

	currentUser := helpers.GetCurrentUser(ctx)
	var user models.User
	dsnap, _ := currentUser.Reference.Get(instance.CtxBackground)
	dsnap.DataTo(&user)

	var posts []IPost.Info
	for userID := range user.Following {
		userRef := usersCollection.Doc(userID)

		dsnap, err := userRef.Get(instance.CtxBackground)
		if err != nil {
			continue
		}

		var author models.User
		dsnap.DataTo(&author)

		docs, _ := postsCollection.Where("author", "==", userRef).OrderBy("createdAt", firestore.Desc).Documents(instance.CtxBackground).GetAll()
		for _, doc := range PostHelpers.Paginate(docs, (query.Page-1)*query.Limit, query.Limit) {
			var post models.Post
			doc.DataTo(&post)

			var likeStatus int
			if post.Likes[currentUser.ID] {
				likeStatus = 1
			} else {
				likeStatus = 0
			}

			var followStatus int
			if helpers.GetCurrentUser(ctx).ID == "-1" || helpers.GetCurrentUser(ctx).ID == author.ID {
				followStatus = -1
			} else if author.Followers[helpers.GetCurrentUser(ctx).ID] {
				followStatus = 1
			} else {
				followStatus = 0
			}

			posts = append(posts, IPost.Info{
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
					ID:           author.ID,
					Name:         author.Name,
					AvatarUrl:    author.AvatarUrl,
					FollowStatus: followStatus,
					Job:          author.Job,
				},
			})
		}
	}

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data:    posts,
	})
}
