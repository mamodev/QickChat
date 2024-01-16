package api

import (
	"github.com/gofiber/fiber/v2"
)


func Logout(c *fiber.Ctx) error {

	c.Accepts("application/json")

	c.ClearCookie("access_token")
	c.ClearCookie("refresh_token")
	
	return c.JSON(fiber.Map{
		"message": "Logged out",
	})
}
