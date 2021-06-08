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

type NewMediaFiles struct {
	ID     string   `json:"id"`
	Images []string `json:"images"`
	Videos []string `json:"videos"`
}

func UploadMediaFiles(ctx iris.Context) {
	postsCollection := services.GetInstance().StoreClient.Collection("posts")

	id := helpers.GetCurrentUserID(ctx)
	refUser := services.GetInstance().StoreClient.Collection("users").Doc(id)

	postID := ctx.Params().Get("id")

	dsnap, err := postsCollection.Doc(postID).Get(instance.CtxBackground)
	if err != nil {
		ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "Post is inexistent"})
		return
	}

	var post models.Post
	dsnap.DataTo(&post)

	if post.Author.ID != refUser.ID {
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

	var images []string
	var videos []string

	for _, file := range files {
		newNameFile := helpers.RandomString(32) + filepath.Ext(file.Filename)
		dest := filepath.Join("./uploads", file.Filename)

		isImage := helpers.IsImage(dest)
		isVideo := helpers.IsVideo(dest)

		if !isImage && !isVideo {
			ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
				Message: "Type of file not supported",
			})
			return
		}

		url, err := helpers.UploadFileStorage(dest, "media/"+post.ID+"/"+newNameFile)
		if err != nil {
			ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
				Message: "Upload media files failed",
			})
			return
		}

		if isImage {
			images = append(images, url)
		} else {
			videos = append(videos, url)
		}
	}

	var listOldFile []string
	if post.Images != nil && post.Videos != nil {
		listOldFile = append(post.Images, post.Videos...)
	} else if post.Images != nil {
		listOldFile = post.Images
	} else {
		listOldFile = post.Videos
	}
	for _, file := range listOldFile {
		helpers.DeleteFileStorage(file)
	}

	_, _ = postsCollection.Doc(postID).Update(instance.CtxBackground, []firestore.Update{
		{
			Path:  "images",
			Value: images,
		},
		{
			Path:  "videos",
			Value: videos,
		},
	})

	helpers.DeleteDir("uploads")

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: NewMediaFiles{
			ID:     post.ID,
			Images: images,
			Videos: videos,
		},
	})
}
