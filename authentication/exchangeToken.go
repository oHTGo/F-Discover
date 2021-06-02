package authentication

import (
	"encoding/json"
	"errors"
	"f-discover/env"
	"f-discover/instance"
	"f-discover/interfaces"
	"f-discover/models"
	"f-discover/services"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

type TokenDTO struct {
	Token string `json:"token"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

var id string
var name string
var avatarUrl string

func ExchangeToken(ctx iris.Context) {
	usersCollection := services.GetInstance().StoreClient.Collection("users")

	var body TokenDTO

	if err := ctx.ReadBody(&body); err != nil {
		ctx.StopWithJSON(iris.StatusUnauthorized, interfaces.IFail{Message: "Token is not a valid token"})
		return
	}

	token := body.Token

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

	j := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Unix() + 60*60*24, //set expire time: 24 hours
	})

	jString, _ := j.SignedString([]byte(env.Get().JWT_SECRET))

	res := TokenResponse{
		Token: jString,
	}

	ctx.JSON(interfaces.ISuccess{
		Message: "Success",
		Data:    res,
	})
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

func checkFieldMapFirebase(payload map[string]interface{}, key string, valueDefault string) string {
	if value, ok := payload[key]; ok {
		return value.(string)
	} else {
		return valueDefault
	}
}
