package api

import (
	"database/sql"
	"log"
	"server/db"

	"github.com/gofiber/fiber/v2"
)

var addFriend = `
	INSERT INTO friend (user1_id, user2_id)
	VALUES ($1, $2), ($2, $1)
`

var deleteFriendRequest = `
	DELETE FROM friend_request
	WHERE id = $1
`

type acceptFriendRequestBody struct {
	Accepted bool `json:"accepted"`
}

func FriendRequestRespond (c *fiber.Ctx) error {
	
	c.Accepts("application/json")

	currentUser := c.Locals("userID").(string)
	requestId := c.Params("id")

	var body acceptFriendRequestBody

	err := c.BodyParser(&body)
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := Validator.Struct(body); err != nil {	
		return fiber.ErrBadRequest
	}

	conn, err := db.Conn(c.Context())
	if err != nil {
		log.Println("Error connecting to database:", err)
		return fiber.ErrInternalServerError
	}
	defer conn.Close()


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

	if receiver_id != currentUser {
		return fiber.ErrForbidden
	}

	if body.Accepted {	
		_, err = tx.ExecContext(c.Context(), addFriend, currentUser, sender_id)
		if err != nil {
			log.Println("Error querying database:", err)
			return fiber.ErrInternalServerError
		}
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
		"message": "Friend request accepted",
	})
}
