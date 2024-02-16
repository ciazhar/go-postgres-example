package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"strings"
)

func GetUserByContext(c *fiber.Ctx) (UserClaims, error) {
	//get authorization header
	authorizationHeader := c.Get("Authorization")
	if authorizationHeader == "" {
		return UserClaims{}, errors.New("Header Not Found")
	}
	//parse header
	headers := strings.Split(authorizationHeader, " ")
	if len(headers) != 2 {
		return UserClaims{}, errors.New("Header Not Valid")
	}

	//get access token
	accessToken := headers[1]

	//get token
	claims := new(CustomClaims)
	token, err := jwt.ParseWithClaims(accessToken, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("jwt.key")), nil
		})
	if err != nil {
		return UserClaims{}, errors.New("Token Not Valid")
	}

	//user
	user := token.Claims.(*CustomClaims).User

	return user, nil
}
