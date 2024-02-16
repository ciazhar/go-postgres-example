package rest

import (
	"github.com/ciazhar/golang-example/generated/db"
	"github.com/ciazhar/golang-example/internal/auth/model"
	"github.com/ciazhar/golang-example/internal/auth/usecase"
	fiber2 "github.com/ciazhar/golang-example/pkg/fiber"
	"github.com/ciazhar/golang-example/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	GetAllRole(c *fiber.Ctx) error
	ForgotPassword(c *fiber.Ctx) error
	CheckPhoneNumber(c *fiber.Ctx) error
}

type authController struct {
	AuthUseCase usecase.AuthUseCase
}

func (it authController) CheckPhoneNumber(c *fiber.Ctx) error {
	id := c.Params("phone")

	exist, err := it.AuthUseCase.CheckPhoneNumber(id)
	if err != nil {
		return err
	}

	return response.Success(c, exist)
}

// ForgotPassword godoc
// @Tags Auth
// @Summary Forgot Password
// @Accept  json
// @Produce  json
// @param data body db.ForgotPasswordParams true "Forgot Password data to be inserted"
// @Success 200 {object} string
// @Router /v1/auth/forgot-password [post]
func (it authController) ForgotPassword(c *fiber.Ctx) error {
	var payload db.ForgotPasswordParams
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(err, response.CodeBadRequest)
	}

	err := it.AuthUseCase.ForgotPassword(payload)
	if err != nil {
		return err
	}

	return response.Success(c, nil)
}

// GetAllRole godoc
// @Tags Auth
// @Summary Data Role
// @Description Only User With This Group Can Access : Admin
// @Accept  json
// @Produce  json
// @Success 200 {object} db.FetchRoleRow
// @Router /v1/auth/role [get]
// @Security Bearer
func (it authController) GetAllRole(c *fiber.Ctx) error {
	res, err := it.AuthUseCase.FetchRole()
	if err != nil {
		return err
	}

	return response.Success(c, res)
}

// Register godoc
// @Tags Auth
// @Summary Register
// @Description All User Can Access This Endpoints, Default Group New User Role is User
// @Accept  json
// @Produce  json
// @param data body db.RegisterParams true "Auth data to be inserted"
// @Success 200 {object} model.TokenResponse
// @Router /v1/auth/register [post]
func (it authController) Register(c *fiber.Ctx) error {
	var payload db.RegisterParams
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(err, response.CodeBadRequest)
	}

	res, err := it.AuthUseCase.Register(payload)
	if err != nil {
		return err
	}

	accessCookie, refreshCookie := fiber2.GetAuthCookies(res.AccessToken, res.RefreshToken)
	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return response.Success(c, res)
}

// Login godoc
// @Tags Auth
// @Summary Login.
// @Description All User Can Access This Endpoints
// @Accept json
// @Produce json
// @param data body model.LoginForm true "Login Form"
// @Success 200 {object} model.TokenResponse
// @Router /v1/auth/login [put]
func (it authController) Login(c *fiber.Ctx) error {
	var payload model.LoginForm
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(err, response.CodeBadRequest)
	}

	res, err := it.AuthUseCase.Login(payload)
	if err != nil {
		return err
	}

	accessCookie, refreshCookie := fiber2.GetAuthCookies(res.AccessToken, res.RefreshToken)
	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return response.Success(c, res)
}

func NewAuthController(AuthUseCase usecase.AuthUseCase) AuthController {
	return authController{AuthUseCase: AuthUseCase}
}
