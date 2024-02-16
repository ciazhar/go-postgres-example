package auth

import (
	"github.com/ciazhar/golang-example/generated/db"
	"github.com/ciazhar/golang-example/internal/auth/controller/rest"
	"github.com/ciazhar/golang-example/internal/auth/usecase"
	usecase2 "github.com/ciazhar/golang-example/internal/token/usecase"
	"github.com/ciazhar/golang-example/internal/user/model"
	"github.com/ciazhar/golang-example/pkg/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Init(router fiber.Router, pool *pgxpool.Pool, db *db.Queries) {
	tokenUC := usecase2.NewTokenUseCase(db, pool)
	uc := usecase.NewAuthUseCase(db, pool, tokenUC)
	controller := rest.NewAuthController(uc)

	r := router.Group("/auth")
	r.Get("/role", auth.ACL(model.Admin), controller.GetAllRole)
	r.Post("/register", controller.Register)
	r.Get("/check-phone/:phone", controller.CheckPhoneNumber)
	r.Post("/forgot-password", controller.ForgotPassword)
	r.Put("/login", controller.Login)
}
