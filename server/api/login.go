package api

import (
	"database/sql"
	"log"
	"server/auth"
	"server/db"

	"github.com/gofiber/fiber/v2"
)

type loginBody struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func Login(c *fiber.Ctx) error {

	c.Accepts("application/json")

	b := new(loginBody)

	if err := c.BodyParser(b); err != nil {
		return fiber.ErrBadRequest
	}

	if err := Validator.Struct(b); err != nil {
		return fiber.ErrBadRequest
	}

	conn, err := db.Conn(c.Context())
	if err != nil {
		log.Println("Error connecting to database:", err)
		return fiber.ErrInternalServerError
	}
	defer conn.Close()

	var password string
	var id string

	err = conn.QueryRowContext(c.Context(), "SELECT id, password FROM users WHERE username = ?", b.Username).Scan(&id, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.ErrUnauthorized
		}

		log.Println("Error querying database:", err)
		return fiber.ErrInternalServerError
	}

	if !auth.ComparePassword(b.Password, password) {
		return fiber.ErrUnauthorized
	}

	accessCookie, refreshCookie, err := auth.CreateTokenCookies(id)
	if err != nil {
		log.Println("Error creating tokens:", err)
		return fiber.ErrInternalServerError
	}

	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return c.JSON(fiber.Map{
		"user_id": id,
		"message": "Successfully logged in",
	})
}
