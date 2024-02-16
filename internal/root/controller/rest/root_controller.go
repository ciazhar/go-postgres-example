package rest

import (
	"github.com/ciazhar/golang-example/internal/root/controller/model"
	"github.com/ciazhar/golang-example/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type RootController interface {
	Root(c *fiber.Ctx) error
}

type rootRestController struct{}

func NewRootRestController() RootController {
	return &rootRestController{}
}

func (r rootRestController) Root(c *fiber.Ctx) error {
	data := model.Data{
		Service: viper.GetString("name"),
		Version: viper.GetString("version"),
		Profile: viper.GetString("profile"),
		Doc:     c.Request().URI().String() + "v1/swagger/index.html",
	}
	return response.Success(c, data)
}
