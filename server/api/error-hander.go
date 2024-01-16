package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler (ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	var e *fiber.Error
	if errors.As(err, &e) {
			code = e.Code
			message = e.Message
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	ctx.Status(code)

	return ctx.JSON(fiber.Map{
			"code":    code,
			"message": message,
	})
}