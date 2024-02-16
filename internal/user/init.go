package user

import (
	"github.com/ciazhar/golang-example/generated/db"
	"github.com/ciazhar/golang-example/internal/user/controller/rest"
	"github.com/ciazhar/golang-example/internal/user/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Init(router fiber.Router, pool *pgxpool.Pool, db *db.Queries) {
	uc := usecase.NewUserUseCase(db, pool)
	controller := rest.NewUserController(uc)

	r := router.Group("/user")
	r.Get("/", controller.GetUser)
	r.Get("/me", controller.Me)
	r.Put("/", controller.UpdateUser)
}
