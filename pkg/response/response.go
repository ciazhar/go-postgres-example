package response

import (
	"errors"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Code    Code        `json:"-"`
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total,omitempty"`
	Count   int         `json:"count,omitempty"`
}

type Code int

const (
	CodeSuccess             Code = 200
	CodeBadRequest               = 400
	CodeInternalServerError      = 500
	CodeUnauthorized             = 401
	CodeForbidden                = 403
)

const SuccessResponse = "success"

func Data(data interface{}, count []int64) Response {
	var countInt int64
	if len(count) != 0 {
		countInt = count[0]
	}
	return Response{
		Count:   int(countInt),
		Message: "succes",
		Code:    CodeSuccess,
		Data:    data,
	}
}

func Success(ctx *fiber.Ctx, data interface{}, count ...int64) error {
	return ctx.JSON(Data(data, count))
}

func Error(err error, code ...Code) error {
	if err == nil {
		return nil
	}

	if code != nil && len(code) == 1 {
		return fiber.NewError(int(code[0]), err.Error())
	} else {
		sentry.CaptureException(err)
		return fiber.NewError(CodeInternalServerError, err.Error())
	}
}

func ErrorS(err string, code ...Code) error {
	if err == "" {
		return nil
	}
	e := errors.New(err)

	if code != nil && len(code) == 1 {
		return fiber.NewError(int(code[0]), e.Error())
	} else {
		sentry.CaptureException(e)
		return fiber.NewError(CodeInternalServerError, e.Error())
	}
}
