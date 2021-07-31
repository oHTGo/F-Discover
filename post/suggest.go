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
	"math/rand"

	"cloud.google.com/go/firestore"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kataras/iris/v12"
)

type SuggestQuery struct {
	Page  int   `url:"page" json:"page"`
	Limit int   `url:"limit" json:"limit"`
	Time  int64 `url:"time" json:"time"`
}

type SuggestResponse struct {
	Posts []IPost.Info `json:"posts"`
}

func Suggest(ctx iris.Context) {
	postsCollection := services.GetInstance().StoreClient.Collection("posts")

	var query SuggestQuery
	if err := ctx.ReadQuery(&query); err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: "Query is invalid",
		})
		return
	}
	if errValidation := validation.ValidateStruct(&query,
		validation.Field(&query.Page, validation.Required, validation.Min(1)),
		validation.Field(&query.Limit, validation.Required, validation.Min(1)),
		validation.Field(&query.Time, validation.Required, validation.Min(1)),
	); errValidation != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFailWithErrors{Message: "Have validation error", Errors: errValidation})
		return
	}

	var posts []IPost.Info

	docs, _ := postsCollection.OrderBy("createdAt", firestore.Desc).Documents(instance.CtxBackground).GetAll()
	var newDocs []*firestore.DocumentSnapshot

	if len(docs) > 0 {
		{
			var randomgen *rand.Rand
			var time int64 = query.Time
			randomgen = rand.New(rand.NewSource(time))
			for _, position := range randomgen.Perm(len(docs) - 1) {
				newDocs = append(newDocs, docs[position])
			}
		}

		for _, doc := range PostHelpers.Paginate(newDocs, (query.Page-1)*query.Limit, query.Limit) {
			var post models.Post
			doc.DataTo(&post)

			dsnap, _ := post.Author.Get(instance.CtxBackground)
			var author models.User
			dsnap.DataTo(&author)

			posts = append(posts, IPost.Info{
				ID:           post.ID,
				Content:      post.Content,
				ThumbnailUrl: post.ThumbnailUrl,
				VideoUrl:     post.VideoUrl,
				Likes:        len(post.Likes),
				LikeStatus:   helpers.GetLikeStatus(ctx, post),
				Comments:     len(post.Comments),
				Location:     location.GetName(post.Location),
				CreatedAt:    post.CreatedAt,
				Author: IPost.Author{
					ID:           author.ID,
					Name:         author.Name,
					AvatarUrl:    author.AvatarUrl,
					FollowStatus: helpers.GetFollowStatus(ctx, author),
					Job:          author.Job,
				},
			})
		}
	}

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: SuggestResponse{
			Posts: posts,
		},
	})
}
