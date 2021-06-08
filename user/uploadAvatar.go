package user

import (
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/services"
	"path/filepath"

	"cloud.google.com/go/firestore"
	"github.com/kataras/iris/v12"
)

type NewAvatarUrl struct {
	AvatarUrl string `json:"avatarUrl"`
}

func UploadAvatar(ctx iris.Context) {
	helpers.CreateDir("uploads")

	id := helpers.GetCurrentUserID(ctx)

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

	newNameFile := id + filepath.Ext(files[0].Filename)
	dest := filepath.Join("./uploads", files[0].Filename)

	if !helpers.IsImage(dest) {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: "Type of file not supported",
		})
		return
	}

	url, err := helpers.UploadFileStorage(dest, "avatar/"+newNameFile)
	if err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: "Upload avatar failed",
		})
		return
	}

	usersCollection := services.GetInstance().StoreClient.Collection("users")
	_, _ = usersCollection.Doc(id).Update(instance.CtxBackground, []firestore.Update{
		{
			Path:  "avatarUrl",
			Value: url,
		},
	})

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data: NewAvatarUrl{
			AvatarUrl: url,
		},
	})
}
