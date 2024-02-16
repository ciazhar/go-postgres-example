package auth

import (
	"errors"
	"github.com/ciazhar/golang-example/pkg/response"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type UserClaims struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	RoleName string `json:"role_name"`
}

type CustomClaims struct {
	jwt.StandardClaims
	User UserClaims `json:"user"`
}

// ACL returns a middleware which secures all the private routes with role
func ACL(role ...string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		//get authorization header
		authorizationHeader := c.Get("Authorization")
		if authorizationHeader == "" {
			return response.Error(errors.New("Header Not Found"), response.CodeUnauthorized)
		}

		//parse header
		headers := strings.Split(authorizationHeader, " ")
		if len(headers) != 2 {
			return response.Error(errors.New("Header Not Valid"), response.CodeUnauthorized)
		}

		//get access token
		accessToken := headers[1]

		//get token
		claims := new(CustomClaims)
		token, err := jwt.ParseWithClaims(accessToken, claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(viper.GetString("jwt.key")), nil
			})

		//validate token
		if token.Valid {
			if claims.ExpiresAt < time.Now().Unix() {
				return response.Error(errors.New("Token Expired"), response.CodeUnauthorized)
			}
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// this is not even a token, we should delete the cookies here
				c.ClearCookie("access_token", "refresh_token")
				return c.SendStatus(fiber.StatusForbidden)
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				return c.SendStatus(fiber.StatusUnauthorized)
			}
		}

		c.Locals("id", claims.Issuer)

		//validate role
		for i := range role {
			if role[i] == claims.User.RoleName {
				return c.Next()
			}
		}
		return response.Error(errors.New("Wrong Role"), response.CodeForbidden)
	}
}
