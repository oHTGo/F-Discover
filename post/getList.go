package post

import (
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	IPost "f-discover/post/interfaces"
	"f-discover/services"

	"github.com/kataras/iris/v12"
	"google.golang.org/api/iterator"
)

func GetListOfUser(ctx iris.Context) {
	postsCollection := services.GetInstance().StoreClient.Collection("posts")
	usersCollection := services.GetInstance().StoreClient.Collection("users")

	userID := ctx.Params().Get("id")
	userRef := usersCollection.Doc(userID)

	if _, err := userRef.Get(instance.CtxBackground); err != nil {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "User is inexistent"})
		return
	}

	var posts []IPost.InfoWithoutAuthor

	iter := postsCollection.Where("author", "==", userRef).Documents(instance.CtxBackground)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			ctx.StopWithJSON(iris.StatusInternalServerError, interfaces.IFail{Message: "Internal Server Error"})
			return
		}
		var post models.Post
		doc.DataTo(&post)

		posts = append(posts, IPost.InfoWithoutAuthor{
			ID:      post.ID,
			Content: post.Content,
			Images:  post.Images,
			Videos:  post.Videos,
			Likes:   len(post.Likes),
		})
	}

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data:    posts,
	})
}
