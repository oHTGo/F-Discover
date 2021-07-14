package post

import (
	"f-discover/env"
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	"f-discover/services"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/kataras/iris/v12"
)

type NewVideo struct {
	ID           string `json:"id"`
	ThumbnailUrl string `json:"thumbnailUrl"`
	VideoUrl     string `json:"videoUrl"`
}

func UploadVideo(ctx iris.Context) {
	postsCollection := services.GetInstance().StoreClient.Collection("posts")

	var postID string = ctx.Params().Get("id")

	var isCreated bool = false
	if postID == "0" {
		isCreated = true
	}

	var post models.Post

	if !isCreated {
		dsnap, err := postsCollection.Doc(postID).Get(instance.CtxBackground)
		if err != nil {
			ctx.StopWithJSON(iris.StatusNotFound, interfaces.IFail{Message: "Post is inexistent"})
			return
		}

		dsnap.DataTo(&post)

		currentUser := helpers.GetCurrentUser(ctx)

		if post.Author.ID != currentUser.ID {
			ctx.StopWithJSON(iris.StatusForbidden, interfaces.IFail{Message: "User is not the author of the post"})
			return
		}
	} else {
		postID = postsCollection.NewDoc().ID
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

	newNameVideo := helpers.RandomString(32) + strings.ToLower(filepath.Ext(files[0].Filename))
	pathLocalVideo := filepath.Join("./uploads", files[0].Filename)
	if !helpers.IsVideo(pathLocalVideo) {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: "Type of file not supported",
		})
		return
	}

	newPathLocalVideo := filepath.Join("./uploads", newNameVideo)

	if env.Get().VIDEO_COMPRESSION {
		if _, err := exec.Command("ffmpeg", "-i", pathLocalVideo,
			"-vcodec", "libx264",
			"-crf", "28",
			newPathLocalVideo).Output(); err != nil {
			helpers.RenameFile(pathLocalVideo, newPathLocalVideo)
		}
	} else {
		helpers.RenameFile(pathLocalVideo, newPathLocalVideo)
	}

	var videoUrl string
	var errUploaded error
	if env.Get().LOCAL_UPLOAD {
		videoUrl, errUploaded = helpers.UploadFileLocal(newPathLocalVideo, "posts/"+postID+"/"+newNameVideo)
	} else {
		videoUrl, errUploaded = helpers.UploadFileStorage(newPathLocalVideo, "posts/"+postID+"/"+newNameVideo)
	}
	if errUploaded != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: "Upload avatar failed",
		})
		return
	}

	thumbnailVideoSecond := int(helpers.GetDurationVideo(newPathLocalVideo) / 2)

	nameThumbnail := helpers.RandomString(32) + ".jpg"
	pathLocalThumbnail := filepath.Join("./uploads", nameThumbnail)
	thumbnailUrl := ""

	if _, err := exec.Command("ffmpeg",
		"-i", newPathLocalVideo,
		"-vframes", "1",
		"-an",
		"-ss", strconv.Itoa(thumbnailVideoSecond),
		pathLocalThumbnail).Output(); err == nil {
		if env.Get().LOCAL_UPLOAD {
			thumbnailUrl, _ = helpers.UploadFileLocal(pathLocalThumbnail, "posts/"+postID+"/"+nameThumbnail)
		} else {
			thumbnailUrl, _ = helpers.UploadFileStorage(pathLocalThumbnail, "posts/"+postID+"/"+nameThumbnail)
		}
	}

	if !isCreated {
		_, _ = postsCollection.Doc(postID).Update(instance.CtxBackground, []firestore.Update{
			{
				Path:  "videoUrl",
				Value: videoUrl,
			},
			{
				Path:  "thumbnailUrl",
				Value: thumbnailUrl,
			},
		})
	} else {
		_, _ = postsCollection.Doc(postID).Set(instance.CtxBackground, map[string]interface{}{
			"id":           postID,
			"author":       helpers.GetCurrentUser(ctx).Reference,
			"thumbnailUrl": thumbnailUrl,
			"videoUrl":     videoUrl,
			"createdAt":    time.Now(),
		})
	}

	// Delete old file in local/storage
	if post.VideoUrl != "" {
		helpers.DeleteFileLocal(post.VideoUrl)
		helpers.DeleteFileStorage(post.VideoUrl)
	}
	if post.ThumbnailUrl != "" {
		helpers.DeleteFileLocal(post.ThumbnailUrl)
		helpers.DeleteFileStorage(post.ThumbnailUrl)
	}

	helpers.DeleteDir("uploads")

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: NewVideo{
			ID:           postID,
			VideoUrl:     videoUrl,
			ThumbnailUrl: thumbnailUrl,
		},
	})
}
