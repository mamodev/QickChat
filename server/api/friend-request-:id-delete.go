package api

import (
	"database/sql"
	"log"
	"server/db"

	"github.com/gofiber/fiber/v2"
)

func FriendRequestDelete (c *fiber.Ctx) error {
	
	c.Accepts("application/json")

	currentUser := c.Locals("userID").(string)
	requestId := c.Params("id")

	if requestId == "" {
		return fiber.ErrBadRequest
	}

	conn, err := db.Conn(c.Context())
	if err != nil {
		log.Println("Error connecting to database:", err)
		return fiber.ErrInternalServerError
	}
	defer conn.Close();


	tx, err := conn.BeginTx(c.Context(), nil)
	if err != nil {
		log.Println("Error beginning transaction:", err)
		return fiber.ErrInternalServerError
	}

	shouldRollback := true
	defer func () {
		if shouldRollback {
			tx.Rollback()
		}
	}()

	var sender_id string
	var receiver_id string

	err = tx.QueryRowContext(c.Context(), `
		SELECT sender_id, receiver_id
		FROM friend_request
		WHERE id = $1
	`, requestId).Scan(&sender_id, &receiver_id)

	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.ErrNotFound
		}

		log.Println("Error querying database:", err)
		return fiber.ErrInternalServerError
	}

	if sender_id != currentUser {
		return fiber.ErrForbidden
	}

	_, err = tx.ExecContext(c.Context(), deleteFriendRequest, requestId)
	if err != nil {
		log.Println("Error querying database:", err)
		return fiber.ErrInternalServerError
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		return fiber.ErrInternalServerError
	}
	shouldRollback = false

	return c.JSON(fiber.Map{
		"message": "Friend request deleted",
	})
}