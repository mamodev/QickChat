package api

import (
	"log"
	"server/db"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)


type friendRequestBody struct {
	UserId string `json:"user_id"`
}

func FriendRequestCreate (c *fiber.Ctx) error {

	c.Accepts("application/json")

	currentUser := c.Locals("userID").(string)

	var body friendRequestBody
	var requestID string = ""

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

	var exists bool

	err = tx.QueryRowContext(c.Context(), `
		SELECT EXISTS (
			SELECT 1
			FROM friend_request
			WHERE sender_id = $1 AND receiver_id = $2
		)
	`, currentUser, body.UserId).Scan(&exists)

	if err != nil {
		log.Println("Error querying database:", err)
		return fiber.ErrInternalServerError
	}

	if exists {
		return fiber.ErrConflict
	}

	// Check if the user is already friends with the other user
	err = tx.QueryRowContext(c.Context(), `
		SELECT EXISTS (
			SELECT 1
			FROM friend
			WHERE user1_id = $1 AND user2_id = $2
			OR user1_id = $2 AND user2_id = $1
		)
	`, currentUser, body.UserId).Scan(&exists)

	if err != nil {
		log.Println("Error querying database:", err)
		return fiber.ErrInternalServerError
	}

	if exists {
		return fiber.ErrConflict
	}

	// check if the user has already sent a friend request to the other user
	err = tx.QueryRowContext(c.Context(), `
		SELECT EXISTS (
			SELECT 1
			FROM friend_request
			WHERE sender_id = $1 AND receiver_id = $2
		)
	`, body.UserId, currentUser).Scan(&exists)

	if err != nil {
		log.Println("Error querying database:", err)
		return fiber.ErrInternalServerError
	}

	if exists {
		// create a friend record
		_, err = tx.ExecContext(c.Context(), `
			INSERT INTO friend (user1_id, user2_id)
			VALUES ($2, $1), ($1, $2)
		`, currentUser, body.UserId)

		if err != nil {
			log.Println("Error querying database:", err)
			return fiber.ErrInternalServerError
		}

		// delete the friend request
		_, err = tx.ExecContext(c.Context(), `
			DELETE FROM friend_request
			WHERE sender_id = $1 AND receiver_id = $2
		`, body.UserId, currentUser)

		if err != nil {
			log.Println("Error querying database:", err)
			return fiber.ErrInternalServerError
		}
	} else {
		requestID = uuid.New().String()
		_, err = tx.ExecContext(c.Context(), `INSERT INTO friend_request (id, sender_id, receiver_id) VALUES ($1, $2, $3)`, requestID, currentUser, body.UserId)

		if err != nil {
			log.Println("Error querying database:", err)
			return fiber.ErrInternalServerError
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		return fiber.ErrInternalServerError
	}
	
	shouldRollback = false

	if requestID == "" {
		return c.JSON(fiber.Map{
			"request_id": nil,
			"friended": true,
			"message": "Friend request sent",
		})
	}

	return c.JSON(fiber.Map{
		"request_id": requestID,
		"friended": false,
		"message": "Friend request sent",
	})
}