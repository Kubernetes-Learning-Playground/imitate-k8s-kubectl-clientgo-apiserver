package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful/v3"
	"net/http"
	"time"
)

// store all user
var userMap = map[string]string{
	"test":    "test",
	"testUser2":    "testUser2",
}

// User 用户
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`

}


func LoginHandler(request *restful.Request, response *restful.Response) {
	switch request.Request.Method {
	case "POST":
		var user User

		err := json.NewDecoder(request.Request.Body).Decode(&user)
		if err != nil {
			fmt.Fprintf(response, "invalid body")
			return
		}

		if userMap[user.Username] == "" || userMap[user.Username] != user.Password {
			fmt.Fprintf(response, "can not authenticate this user")
			return
		}

		token, err := generateJWT(user.Username)
		if err != nil {
			fmt.Fprintf(response, "error in generating token")
		}

		fmt.Fprintf(response, token)

	}
}

var sampleSecretKey = []byte("api-server-secret-key")

// generateJWT 生成token
func generateJWT(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

	tokenString, err := token.SignedString(sampleSecretKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

// ValidateToken 认证token
func ValidateToken(w http.ResponseWriter, r *http.Request) (err error) {
	if r.Header["Token"] == nil {
		fmt.Fprintf(w, "can not find token in header")
		return errors.New("Token error")
	}

	token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return sampleSecretKey, nil
	})

	if token == nil {
		fmt.Fprintf(w, "invalid token")
		return errors.New("Token error")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Fprintf(w, "couldn't parse claims")
		return errors.New("Token error")
	}

	// 过期时间token
	exp := claims["exp"].(float64)
	if int64(exp) < time.Now().Local().Unix() {
		fmt.Fprintf(w, "token expired")
		return errors.New("Token error")
	}

	// 把username放入header中
	userName := claims["username"].(string)
	r.Header.Set("username", userName)

	return nil
}
