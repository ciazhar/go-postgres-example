package main

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	_ "github.com/ciazhar/golang-example/generated/docs"
	"github.com/ciazhar/golang-example/internal/auth"
	"github.com/ciazhar/golang-example/internal/root"
	"github.com/ciazhar/golang-example/internal/user"
	"github.com/ciazhar/golang-example/pkg/env"
	error_handler "github.com/ciazhar/golang-example/pkg/error"
	logger "github.com/ciazhar/golang-example/pkg/log"
	"github.com/ciazhar/golang-example/pkg/postgres"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
)

// @title iLikes API Documentation
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @BasePath /
func main() {

	//init fiber and middleware
	r := fiber.New(fiber.Config{
		ErrorHandler: error_handler.InitErrorHandler(),
	})
	r.Use(cors.New(cors.Config{AllowCredentials: true}))
	r.Use(recover.New())
	g := r.Group("/v1/")

	env.Init()
	logger.Init()
	pool, dbx := postgres.Init()

	//route internal
	root.Init(r)
	auth.Init(g, pool, dbx)
	user.Init(g, pool, dbx)

	//route swagger
	g.Get("/swagger/*", swagger.HandlerDefault) // default

	//route 404
	r.Use(func(c *fiber.Ctx) error {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "route not found",
		})
	})

	//run
	if viper.GetString("profile") != "debug" {
		sentry.CaptureMessage("appplication start in port : " + viper.GetString("port"))
	}
	err := r.Listen(":" + viper.GetString("port"))
	if err != nil {
		panic(err)
	}
}
