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

	id := ctx.Values().GetString("id")

	_, fileHeader, err := ctx.FormFile("avatar")
	if err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: "Upload avatar failed",
		})
		return
	}

	// Upload the file to specific destination.
	newNameFile := id + filepath.Ext(fileHeader.Filename)
	dest := filepath.Join("./uploads", newNameFile)

	_, err = ctx.SaveFormFile(fileHeader, dest)
	if err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: "Upload avatar failed",
		})
		return
	}

	url, err := helpers.UploadFile(dest, "avatar/"+newNameFile)
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
		Message: "Upload avatar successfully",
		Data: NewAvatarUrl{
			AvatarUrl: url,
		},
	})
}
