package root

import (
	"github.com/ciazhar/golang-example/internal/root/controller/rest"
	"github.com/gofiber/fiber/v2"
)

func Init(router fiber.Router) {

	controller := rest.NewRootRestController()

	r := router.Group("/")
	r.Get("/", controller.Root)

}
