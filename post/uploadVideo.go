package post

import (
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	"f-discover/services"
	"path/filepath"

	"cloud.google.com/go/firestore"
	"github.com/kataras/iris/v12"
)

type NewVideo struct {
	ID       string `json:"id"`
	VideoUrl string `json:"videoUrl"`
}

func UploadVideo(ctx iris.Context) {
	postsCollection := services.GetInstance().StoreClient.Collection("posts")

	postID := ctx.Params().Get("id")

	dsnap, err := postsCollection.Doc(postID).Get(instance.CtxBackground)
	if err != nil {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "Post is inexistent"})
		return
	}

	var post models.Post
	dsnap.DataTo(&post)

	currentUser := helpers.GetCurrentUser(ctx)

	if post.Author.ID != currentUser.ID {
		ctx.StopWithJSON(iris.StatusForbidden, interfaces.IFail{Message: "User is not the author of the post"})
		return
	}

	files, err := helpers.UploadFiles(ctx)
	if err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: err.Error(),
		})
		return
	}

	if len(files) != 1 {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: "Only upload 1 file",
		})
		return
	}

	newNameFile := postID + "/" + helpers.RandomString(32) + filepath.Ext(files[0].Filename)
	dest := filepath.Join("./uploads", files[0].Filename)

	if !helpers.IsVideo(dest) {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: "Type of file not supported",
		})
		return
	}

	url, err := helpers.UploadFileStorage(dest, "posts/"+newNameFile)
	if err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: "Upload avatar failed",
		})
		return
	}

	_, _ = postsCollection.Doc(postID).Update(instance.CtxBackground, []firestore.Update{
		{
			Path:  "videoUrl",
			Value: url,
		},
	})

	if post.VideoUrl != "" {
		helpers.DeleteFileStorage(post.VideoUrl)
	}
	helpers.DeleteDir("uploads")

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: NewVideo{
			ID:       post.ID,
			VideoUrl: url,
		},
	})
}
