package api

import (
	"log"
	"server/db"
	"server/events"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ChatMessageGet (c *fiber.Ctx) error {

	c.Accepts("application/json")

	chatID := c.Params("id")
	userID := c.Locals("userID").(string)

	unsafeBeforeTime := c.Query("before_time")
	beforeTime := time.Now()

	if unsafeBeforeTime != "" {
		var err error
		beforeTime, err = time.Parse(time.RFC3339, unsafeBeforeTime)
		if err != nil {
			return fiber.ErrBadRequest
		}
	}

	conn, err := db.Conn(c.Context())
	if err != nil {
		log.Println("Error connecting to database:", err)
		return fiber.ErrInternalServerError
	}

	defer conn.Close()


	// check if user is in chat
	var chatExists bool
	err = conn.QueryRowContext(c.Context(), `
		SELECT EXISTS (
			SELECT 1
			FROM chat_user
			WHERE chat_id = $1 AND user_id = $2
		)
	`, chatID, userID).Scan(&chatExists)

	if err != nil {
		log.Println("Error querying database:", err)
		return fiber.ErrInternalServerError
	}

	if !chatExists {
		return fiber.ErrForbidden
	}

	sqlTime := beforeTime.Format("2006-01-02 15:04:05")
	rows, err := conn.QueryContext(c.Context(), `
		SELECT
			message.id,
			chat_id,
			message,
			sender_id,
			username,
			message.created_at
		FROM message join users on sender_id = users.id
		WHERE message.created_at < $1 AND chat_id = $2
		ORDER BY message.created_at DESC
		LIMIT 50
	`, sqlTime, chatID)
	
	if err != nil {
		log.Println("Error querying database:", err)
		return fiber.ErrInternalServerError
	}

	defer rows.Close()

	messages := []events.Message{}

	for rows.Next() {
		var message events.Message
		err := rows.Scan(
			&message.ID,
			&message.ChatID,
			&message.Message,
			&message.SenderID,
			&message.SenderUsername,
			&message.Timestamp,
		)

		if err != nil {
			log.Println("Error scanning row:", err)
			return fiber.ErrInternalServerError
		}

		messages = append(messages, message)
	}

	return c.JSON(messages)
}		

