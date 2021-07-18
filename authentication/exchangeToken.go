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
	"strings"
	"time"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

type TokenDTO struct {
	Type  string `json:"type"`
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
		ctx.StopWithJSON(iris.StatusUnauthorized, interfaces.IFail{Message: "Body is bad request"})
		return
	}

	var user *User

	switch strings.ToLower(body.Type) {
	case "firebase":
		if decoded, err := verifyTokenFirebase(body.Token); err == nil {
			user = decoded
		} else {
			ctx.StopWithJSON(iris.StatusUnauthorized, interfaces.IFail{Message: err.Error()})
			return
		}
	case "zalo":
		token, err := exchangeTokenZalo(body.Token)
		if err != nil {
			ctx.StopWithJSON(iris.StatusUnauthorized, interfaces.IFail{Message: err.Error()})
			return
		}

		if decoded, err := verifyTokenZalo(token); err == nil {
			user = decoded
		} else {
			ctx.StopWithJSON(iris.StatusUnauthorized, interfaces.IFail{Message: "Token not verified"})
			return
		}
	default:
		ctx.StopWithJSON(iris.StatusUnauthorized, interfaces.IFail{Message: "Type is invalid"})
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
		"exp": time.Now().Unix() + 60*60*24*30, //set expire time: 30 days
	})

	jString, _ := j.SignedString([]byte(env.Get().JWT_SECRET))
	ctx.SetCookieKV("token", jString, iris.CookieCleanPath, iris.CookieHTTPOnly(false))

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
		return nil, errors.New("Token not verified")
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

func exchangeTokenZalo(code string) (string, error) {
	APP_ID_ZALO := env.Get().APP_ID_ZALO
	APP_SECRET_ZALO := env.Get().APP_SECRET_ZALO
	url := "https://oauth.zaloapp.com/v3/access_token?app_id=" + APP_ID_ZALO + "&app_secret=" + APP_SECRET_ZALO + "&code=" + code
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New("Code is invalid")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Code is invalid")
	}

	mBody := make(map[string]interface{})
	if err = json.Unmarshal(body, &mBody); err != nil {
		return "", errors.New("Code is invalid")
	}

	if _, ok := mBody["access_token"]; ok {
		return mBody["access_token"].(string), nil
	} else {
		return "", errors.New("Code is invalid")
	}
}

func verifyTokenZalo(token string) (*User, error) {
	url := "https://graph.zalo.me/v2.0/me?access_token=" + token + "&fields=id,name,picture"
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("Token not verified")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Token not verified")
	}

	mBody := make(map[string]interface{})
	if err = json.Unmarshal(body, &mBody); err != nil {
		return nil, errors.New("Code is invalid")
	}

	if errorCode, ok := mBody["error"]; ok {
		if errorCode.(float64) != 0 {
			return nil, errors.New("Token not verified")
		}
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
