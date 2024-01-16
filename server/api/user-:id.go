package api

import (
	"database/sql"
	"log"
	"server/db"

	"github.com/gofiber/fiber/v2"
)

func UserGet(c *fiber.Ctx) error {

	userID := c.Params("id")

	conn, err := db.Conn(c.Context())
	if err != nil {
		log.Println("Error connecting to database:", err)
		return fiber.ErrInternalServerError
	}
	defer conn.Close()
	
	var user struct {
		Id string `json:"id"`
		Username string `json:"username"`
		Email string `json:"email"`
	}
	err = conn.QueryRowContext(c.Context(), `
		select
			id,
			username,
			email
		from users
		where id = $1
	`, userID).Scan(&user.Id, &user.Username, &user.Email)

	if err != nil {
		if sql.ErrNoRows == err {
			return fiber.ErrNotFound
		}

		log.Println("Error querying database:", err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(user)
}