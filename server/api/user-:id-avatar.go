package api

import (
	"database/sql"
	"log"
	"server/db"

	"github.com/gofiber/fiber/v2"
)

func UserAvatarGet (c *fiber.Ctx) error {
	c.Accepts("image/png", "image/jpeg")

	id := c.Params("id")

	conn, err := db.Conn(c.Context())
	if err != nil {
		log.Println("Error connecting to database:", err)
		return fiber.ErrInternalServerError
	}
	defer conn.Close()
	

	var path string
	err = conn.QueryRowContext(c.Context(), `
		select
			profile_picture
		from users
		where id = $1
	`, id).Scan(&path)

	if err != nil {
		if sql.ErrNoRows == err {
			return fiber.ErrNotFound
		}

		log.Println("Error querying database:", err)
		return fiber.ErrInternalServerError
	}

	return c.SendFile("../database/img/" + path)
}