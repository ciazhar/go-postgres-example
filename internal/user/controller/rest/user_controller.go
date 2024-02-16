package rest

import (
	"github.com/ciazhar/golang-example/generated/db"
	"github.com/ciazhar/golang-example/internal/user/model"
	"github.com/ciazhar/golang-example/internal/user/usecase"
	"github.com/ciazhar/golang-example/pkg/auth"
	"github.com/ciazhar/golang-example/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserController interface {
	Me(c *fiber.Ctx) error
	GetUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
}

type userController struct {
	UserUseCase usecase.UserUseCase
}

// Me godoc
// @Tags User
// @Summary Fetch user
// @Description Fetch user
// @Accept  json
// @Produce  json
// @Success 200 {object} db.GetUserByIdRow
// @Router /v1/user/me [get]
// @Security Bearer
func (it userController) Me(c *fiber.Ctx) error {
	user, err := auth.GetUserByContext(c)
	if err != nil {
		return response.Error(err, response.CodeUnauthorized)
	}

	res, err := it.UserUseCase.GetUserByID(uuid.MustParse(user.Id))
	if err != nil {
		return response.Error(err, response.CodeInternalServerError)
	}
	return response.Success(c, res)
}

// GetUser godoc
// @Tags User
// @Summary Fetch user
// @Description Fetch user
// @Accept  json
// @Produce  json
// @param page query int true "Page number"
// @param size query int true "Size row"
// @Success 200 {object} db.FetchUserRow
// @Router /v1/user [get]
// @Security Bearer
func (it userController) GetUser(c *fiber.Ctx) error {
	param := model.FetchParam{}
	if err := c.QueryParser(&param); err != nil {
		return response.Error(err, response.CodeBadRequest)
	}

	payload, err := it.UserUseCase.GetUser(param)
	if err != nil {
		return err
	}

	return response.Success(c, payload)
}

// UpdateUser godoc
// @Tags User
// @Summary Update User
// @Description Update User
// @Accept json
// @Produce json
// @param data body db.UpdateUserParams true "Update User Form"
// @Router /v1/user [put]
func (it userController) UpdateUser(c *fiber.Ctx) error {
	var payload db.UpdateUserParams
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(err, response.CodeBadRequest)
	}

	err := it.UserUseCase.UpdateUser(payload)
	if err != nil {
		return err
	}

	return response.Success(c, nil)
}
func NewUserController(UserUseCase usecase.UserUseCase) UserController {
	return userController{UserUseCase: UserUseCase}
}
