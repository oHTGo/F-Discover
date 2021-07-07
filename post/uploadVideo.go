package post

import (
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	"f-discover/services"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

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
	//postID + "/" +
	newNameVideo := helpers.RandomString(32) + strings.ToLower(filepath.Ext(files[0].Filename))
	pathLocalVideo := filepath.Join("./uploads", files[0].Filename)

	if !helpers.IsVideo(pathLocalVideo) {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: "Type of file not supported",
		})
		return
	}

	videoUrl, err := helpers.UploadFileStorage(pathLocalVideo, "posts/"+postID+"/"+newNameVideo)
	if err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: "Upload video failed",
		})
		return
	}

	newPathLocalVideo := filepath.Join("./uploads", helpers.RandomString(32)+strings.ToLower(filepath.Ext(files[0].Filename)))
	helpers.RenameFile(pathLocalVideo, newPathLocalVideo)

	secondToGenerateThumbnail := int(helpers.GetDurationVideo(newPathLocalVideo) / 2)

	nameThumbnail := helpers.RandomString(32) + ".jpg"
	pathLocalThumbnail := filepath.Join("./uploads", nameThumbnail)
	thumbnailUrl := ""

	_, errExec := exec.Command("ffmpeg", "-i", newPathLocalVideo, "-vframes", "1", "-an", "-ss", strconv.Itoa(secondToGenerateThumbnail), pathLocalThumbnail).Output()
	if errExec == nil {
		thumbnailUrl, _ = helpers.UploadFileStorage(pathLocalThumbnail, "posts/"+postID+"/"+nameThumbnail)
	}

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

	// Delete old file in storage
	if post.VideoUrl != "" {
		helpers.DeleteFileStorage(post.VideoUrl)
	}
	if post.ThumbnailUrl != "" {
		helpers.DeleteFileStorage(post.ThumbnailUrl)
	}

	helpers.DeleteDir("uploads")

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: NewVideo{
			ID:           post.ID,
			VideoUrl:     videoUrl,
			ThumbnailUrl: thumbnailUrl,
		},
	})
}
