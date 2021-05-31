package middlewares

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	"f-discover/services"

	"github.com/kataras/iris/v12"
)

var id string
var name string
var avatarUrl string

func SetAuthentication() iris.Handler {
	return func(ctx iris.Context) {
		usersCollection := services.GetInstance().StoreClient.Collection("users")

		authorization := ctx.GetHeader("authorization")
		typeOfToken := strings.Split(authorization, " ")[0]

		if typeOfToken == "" || typeOfToken != "Bearer" {
			ctx.StopWithJSON(iris.StatusUnauthorized, interfaces.IFail{Message: "Token is not a valid token"})
			return
		}

		token := strings.Split(authorization, " ")[1]
		if token == "" {
			ctx.StopWithJSON(iris.StatusUnauthorized, interfaces.IFail{Message: "Token is not a valid token"})
			return
		}

		idFirebase, nameFirebase, avatarUrlFirebase, errFirebase := verifyTokenFirebase(token)
		idZalo, nameZalo, avatarUrlZalo, errZalo := verifyTokenZalo(token)

		if errFirebase == nil {
			id = idFirebase
			name = nameFirebase
			avatarUrl = avatarUrlFirebase
		} else if errZalo == nil {
			id = idZalo
			name = nameZalo
			avatarUrl = avatarUrlZalo
		} else {
			ctx.StopWithJSON(iris.StatusUnauthorized, interfaces.IFail{Message: "Token not verified"})
			return
		}

		_, err := usersCollection.Doc(id).Get(instance.CtxBackground)
		if err != nil {
			user := models.User{
				ID:        id,
				Name:      name,
				AvatarUrl: avatarUrl,
			}
			usersCollection.Doc(id).Set(instance.CtxBackground, user)
		}

		ctx.Values().Set("id", id)

		ctx.Next()
	}
}

func checkFieldMapFirebase(payload map[string]interface{}, key string, valueDefault string) string {
	if value, ok := payload[key]; ok {
		return value.(string)
	} else {
		return valueDefault
	}
}

func verifyTokenFirebase(token string) (string, string, string, error) {
	authClient := services.GetInstance().AuthClient
	payload, err := authClient.VerifyIDToken(instance.CtxBackground, token)
	if err != nil {
		return "", "", "", err
	}

	id := payload.Claims["user_id"].(string)
	name := checkFieldMapFirebase(payload.Claims, "name", "")
	avatarUrl := checkFieldMapFirebase(payload.Claims, "picture", "")

	return id, name, avatarUrl, nil
}

func verifyTokenZalo(token string) (string, string, string, error) {
	url := "https://graph.zalo.me/v2.0/me?access_token=" + token + "&fields=id,name,picture"
	resp, err := http.Get(url)
	if err != nil {
		return "", "", "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}

	mBody := make(map[string]interface{})
	json.Unmarshal(body, &mBody)

	if _, ok := mBody["error"]; ok {
		return "", "", "", errors.New("token not verified")
	}

	id := mBody["id"].(string)
	name := mBody["name"].(string)
	avatarUrl := mBody["picture"].(map[string]interface{})["data"].(map[string]interface{})["url"].(string)

	return id, name, avatarUrl, nil
}
