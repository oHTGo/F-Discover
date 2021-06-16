package search

import (
	"f-discover/helpers"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/location"
	"f-discover/models"
	"f-discover/services"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kataras/iris/v12"
)

type Result struct {
	ID   string `json:"ID"`
	Name string `json:"name"`
}

type SearchResponse struct {
	Type   string `json:"type"`
	Result Result `json:"result"`
}

type SearchQuery struct {
	Search string `json:"search" url:"search"`
}

func Search(ctx iris.Context) {
	var query SearchQuery
	if err := ctx.ReadQuery(&query); err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFail{
			Message: "Query is invalid",
		})
		return
	}
	if errValidation := validation.ValidateStruct(&query,
		validation.Field(&query.Search, validation.Required),
	); errValidation != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, interfaces.IFailWithErrors{Message: "Have validation error", Errors: errValidation})
		return
	}

	var res []SearchResponse

	usersCollection := services.GetInstance().StoreClient.Collection("users")
	docs, _ := usersCollection.Documents(instance.CtxBackground).GetAll()
	for _, doc := range docs {
		var user models.User
		doc.DataTo(&user)

		nameLower := strings.ToLower(user.Name)
		searchLower := strings.ToLower(query.Search)
		if strings.Contains(
			nameLower,
			searchLower,
		) || strings.Contains(
			helpers.ConvertUnicodeToASCII(nameLower),
			helpers.ConvertUnicodeToASCII(searchLower),
		) {
			res = append(res, SearchResponse{
				Type: "user",
				Result: Result{
					ID:   user.ID,
					Name: user.Name,
				},
			})
		}
	}

	for _, location := range location.FindByName(query.Search) {
		res = append(res, SearchResponse{
			Type: "location",
			Result: Result{
				ID:   location.ID,
				Name: location.Name,
			},
		})
	}

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data:    res,
	})
}
