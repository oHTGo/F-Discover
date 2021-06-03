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

type User struct {
	ID        string
	Name      string
	AvatarUrl string
}

func ExchangeToken(ctx iris.Context) {
	usersCollection := services.GetInstance().StoreClient.Collection("users")

	var body TokenDTO

	if err := ctx.ReadBody(&body); err != nil {
		ctx.StopWithJSON(iris.StatusUnauthorized, interfaces.IFail{Message: "Token is not a valid token"})
		return
	}

	token := body.Token
	var user *User

	if decoded, err := verifyTokenFirebase(token); err == nil {
		user = decoded
	} else if decoded, err := verifyTokenZalo(token); err == nil {
		user = decoded
	} else {
		ctx.StopWithJSON(iris.StatusUnauthorized, interfaces.IFail{Message: "Token not verified"})
		return
	}

	_, err := usersCollection.Doc(user.ID).Get(instance.CtxBackground)
	if err != nil {
		data := models.User{
			ID:        user.ID,
			Name:      user.Name,
			AvatarUrl: user.AvatarUrl,
		}
		usersCollection.Doc(user.ID).Set(instance.CtxBackground, data)
	}

	j := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
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

func verifyTokenFirebase(token string) (*User, error) {
	authClient := services.GetInstance().AuthClient
	payload, err := authClient.VerifyIDToken(instance.CtxBackground, token)
	if err != nil {
		return nil, err
	}

	var id, name, avatarUrl string

	id = payload.Claims["user_id"].(string)
	if namePayload, ok := payload.Claims["name"]; ok {
		name = namePayload.(string)
	} else {
		name = payload.Claims["phone_number"].(string)
	}
	if avatarUrlPayload, ok := payload.Claims["picture"]; ok {
		avatarUrl = avatarUrlPayload.(string)
	}

	return &User{
		ID:        id,
		Name:      name,
		AvatarUrl: avatarUrl,
	}, nil
}

func verifyTokenZalo(token string) (*User, error) {
	url := "https://graph.zalo.me/v2.0/me?access_token=" + token + "&fields=id,name,picture"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	mBody := make(map[string]interface{})
	json.Unmarshal(body, &mBody)

	if _, ok := mBody["error"]; ok {
		return nil, errors.New("token not verified")
	}

	id := mBody["id"].(string)
	name := mBody["name"].(string)
	avatarUrl := mBody["picture"].(map[string]interface{})["data"].(map[string]interface{})["url"].(string)

	return &User{
		ID:        id,
		Name:      name,
		AvatarUrl: avatarUrl,
	}, nil
}
