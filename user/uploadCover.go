package user

import (
	"f-discover/env"
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/services"
	"log"
	"path/filepath"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/kataras/iris/v12"
)

type NewCoverUrl struct {
	CoverUrl string `json:"coverUrl"`
}

func UploadCover(ctx iris.Context) {
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

	currentUser := helpers.GetCurrentUser(ctx)

	newNameFile := currentUser.ID + strings.ToLower(filepath.Ext(files[0].Filename))
	dest := filepath.Join("./uploads", files[0].Filename)

	if !helpers.IsImage(dest) {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: "Type of file not supported",
		})
		return
	}

	var url string
	var errUploaded error
	if env.Get().LOCAL_UPLOAD {
		log.Println("Uploading")
		url, errUploaded = helpers.UploadFileLocal(dest, "cover/"+newNameFile)
	} else {
		url, errUploaded = helpers.UploadFileStorage(dest, "cover/"+newNameFile)
	}
	if errUploaded != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: "Upload avatar failed",
		})
		return
	}

	usersCollection := services.GetInstance().StoreClient.Collection("users")
	_, _ = usersCollection.Doc(currentUser.ID).Update(instance.CtxBackground, []firestore.Update{
		{
			Path:  "coverUrl",
			Value: url,
		},
	})

	helpers.DeleteDir("uploads")

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: NewCoverUrl{
			CoverUrl: url,
		},
	})
}
