package api

import (
	"server/events"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func ChatUpdate(c *fiber.Ctx) error {

	c.Accepts("application/json")

	userID := utils.CopyString(c.Locals("userID").(string))

	messages := events.Pool.GetUsersMessages(userID)

	return c.JSON(messages)
}