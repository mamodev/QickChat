package api

import (
	"log"
	"server/auth"
	"server/db"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type registerBody struct {
	Username string `json:"username" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func Register(c *fiber.Ctx) (error) {

	c.Accepts("application/json")

	b := new(registerBody)

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

	rows, err := conn.QueryContext(c.Context(), "SELECT id FROM users WHERE username = ? OR email = ?", b.Username, b.Email);
	defer rows.Close()

	for rows.Next() {
		return fiber.NewError(fiber.StatusConflict, "Username or email already exists")
	}

	encrypedPassword, err := auth.EncryptPassword(b.Password)
	if err != nil {
		log.Println("Error encrypting password:", err)
		return fiber.ErrInternalServerError
	}

	uuid := uuid.New().String()
	_, err = conn.ExecContext(c.Context(), "INSERT INTO users (id, username, email, password) VALUES (?, ?, ?, ?)", uuid, b.Username, b.Email, encrypedPassword)

	if err != nil {
		log.Println("Error inserting user into database:", err)
		return fiber.ErrInternalServerError
	}

	accessCookie, refreshCookie, err := auth.CreateTokenCookies(uuid)

	if err != nil {
		log.Println("Error creating tokens:", err)
		return fiber.ErrInternalServerError
	}

	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return c.JSON(fiber.Map{
		"user_id": uuid,
		"message": "User created",
	})
}
