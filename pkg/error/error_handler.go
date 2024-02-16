package error_handler

import (
	"github.com/ciazhar/golang-example/pkg/response"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
)

func InitErrorHandler() func(c *fiber.Ctx, err error) error {
	return func(c *fiber.Ctx, err error) error {
		// Default unknown error
		code := response.CodeInternalServerError

		if e, ok := err.(*fiber.Error); ok {
			// Override code if fiber.Error type
			code = e.Code
		} else {
			// log error
			sentry.CaptureException(err)
		}

		// Set Content-Type: text/plain; charset=utf-8
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Return statuscode with error message
		return c.Status(code).JSON(response.Response{
			Message: "error",
			Error:   err.Error(),
		})
	}
}
