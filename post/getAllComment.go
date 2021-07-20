package post

import (
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	IPost "f-discover/post/interfaces"
	"f-discover/services"
	"sort"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kataras/iris/v12"
)

type GetAllCommentQuery struct {
	Page  int `url:"page" json:"page"`
	Limit int `url:"limit" json:"limit"`
}

type GetAllCommentResponse struct {
	Total    int             `json:"total"`
	Comments []IPost.Comment `json:"comments"`
}

func GetAllComment(ctx iris.Context) {
	postsCollection := services.GetInstance().StoreClient.Collection("posts")

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

	postID := ctx.Params().Get("id")
	dsnap, err := postsCollection.Doc(postID).Get(instance.CtxBackground)
	if err != nil {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "Post is inexistent"})
		return
	}
	var post models.Post
	dsnap.DataTo(&post)

	var sortedComments []models.Comment
	for _, comment := range post.Comments {
		sortedComments = append(sortedComments, models.Comment{
			ID:        comment.ID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
			Author:    comment.Author,
		})
	}
	// Sort by descending creation time
	sort.Slice(sortedComments, func(i, j int) bool {
		return sortedComments[i].CreatedAt.After(sortedComments[j].CreatedAt)
	})

	var comments []IPost.Comment

	for _, comment := range paginate(sortedComments, (query.Page-1)*query.Limit, query.Limit) {
		dsnap, _ = comment.Author.Get(instance.CtxBackground)
		var author models.User
		dsnap.DataTo(&author)

		var followStatus int
		if helpers.GetCurrentUser(ctx).ID == "-1" || helpers.GetCurrentUser(ctx).ID == author.ID {
			followStatus = -1
		} else if author.Followers[helpers.GetCurrentUser(ctx).ID] {
			followStatus = 1
		} else {
			followStatus = 0
		}

		comments = append(comments, IPost.Comment{
			ID:        comment.ID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
			Author: IPost.Author{
				ID:           author.ID,
				Name:         author.Name,
				AvatarUrl:    author.AvatarUrl,
				FollowStatus: followStatus,
				Job:          author.Job,
			},
		})
	}

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: GetAllCommentResponse{
			Total:    len(comments),
			Comments: comments,
		},
	})
}

func paginate(x []models.Comment, skip int, size int) []models.Comment {
	if skip > len(x) {
		skip = len(x)
	}

	end := skip + size
	if end > len(x) {
		end = len(x)
	}

	return x[skip:end]
}
