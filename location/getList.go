package location

import (
	"f-discover/interfaces"

	"github.com/kataras/iris/v12"
)

type GetListResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetList(ctx iris.Context) {
	var list []GetListResponse
	for id, name := range LOCATIONS {
		list = append(list, GetListResponse{
			ID:   id,
			Name: name,
		})
	}

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data:    list,
	})
}
