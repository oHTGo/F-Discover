package post

import (
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

type GetListOfUserQuery struct {
	Page  int `url:"page" json:"page"`
	Limit int `url:"limit" json:"limit"`
}

type GetListOfUserResponse struct {
	Total int                       `json:"total"`
	Posts []IPost.InfoWithoutAuthor `json:"posts"`
}

func GetListOfUser(ctx iris.Context) {
	postsCollection := services.GetInstance().StoreClient.Collection("posts")
	usersCollection := services.GetInstance().StoreClient.Collection("users")

	var query GetListOfUserQuery
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

	userID := ctx.Params().Get("id")
	userRef := usersCollection.Doc(userID)

	if _, err := userRef.Get(instance.CtxBackground); err != nil {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "User is inexistent"})
		return
	}

	var posts []IPost.InfoWithoutAuthor

	docs, _ := postsCollection.Where("author", "==", userRef).OrderBy("createdAt", firestore.Desc).Documents(instance.CtxBackground).GetAll()
	for _, doc := range PostHelpers.Paginate(docs, (query.Page-1)*query.Limit, query.Limit) {
		var post models.Post
		doc.DataTo(&post)

		posts = append(posts, IPost.InfoWithoutAuthor{
			ID:           post.ID,
			Content:      post.Content,
			ThumbnailUrl: post.ThumbnailUrl,
			VideoUrl:     post.VideoUrl,
			Likes:        len(post.Likes),
			Comments:     len(post.Comments),
			Location:     location.GetName(post.Location),
			CreatedAt:    post.CreatedAt,
		})
	}

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: GetListOfUserResponse{
			Total: len(docs),
			Posts: posts,
		},
	})
}
