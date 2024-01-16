package api

import (
	"database/sql"
	"log"
	"server/db"
	"server/events"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
)

type chatMessageCreateRequest struct {
	Message string `json:"message" validate:"required"`
}

func ChatMessageCreate (ctx *fiber.Ctx) error {

	ctx.Accepts("application/json")

	userID := utils.CopyString(ctx.Locals("userID").(string))
	chatID := utils.CopyString(ctx.Params("id"))

	var body chatMessageCreateRequest
	if err := ctx.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := Validator.Struct(body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	
	conn, err := db.Conn(ctx.Context())
	if err != nil {
		log.Println("Error connecting to database:", err)
		return fiber.ErrInternalServerError
	}

	defer conn.Close()

	tx, err := conn.BeginTx(ctx.Context(), nil)
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

	usersID := []string{}
	rows, err := tx.QueryContext(ctx.Context(), `
		SELECT user_id
		FROM chat_user
		WHERE chat_id = $1
	`, chatID)

	if err != nil {
		log.Println("Error querying for chat membership:", err)
		return fiber.ErrInternalServerError
	}

	isMember := false
	for rows.Next() {
		var userID string
		err := rows.Scan(&userID)
		if err != nil {
			log.Println("Error scanning chat membership:", err)
			return fiber.ErrInternalServerError
		}

		if userID == userID {
			isMember = true
		}

		usersID = append(usersID, userID)
	}

	if len(usersID) == 0 {
		return fiber.ErrNotFound
	}

	if !isMember {
		return fiber.ErrForbidden
	}

	msgID := uuid.New().String()


	var username string
	err = tx.QueryRowContext(ctx.Context(), `
		SELECT username
		FROM users
		WHERE id = $1
	`, userID).Scan(&username)

	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.ErrForbidden
		}

		log.Println("Error querying for user:", err)
		return fiber.ErrInternalServerError
	}

	msg := events.Message{
		ID: msgID,
		ChatID: chatID,
		SenderID: userID,
		SenderUsername: username,
		Message: body.Message,
		Timestamp: time.Now(),
	}

	events.Pool.AddMessage(usersID, msg)
	// for _, userID := range usersID {
	// 	// SendNotification(ctx, "You have a new message", userID)
	// }


	_, err = tx.ExecContext(ctx.Context(), `
		INSERT INTO message (id, chat_id, sender_id, message)
		VALUES ($1, $2, $3, $4)
	`, msgID, chatID, userID, body.Message)

	if err != nil {
		log.Println("Error inserting chat message:", err)
		return fiber.ErrInternalServerError
	}


	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		return fiber.ErrInternalServerError
	}

	shouldRollback = false
	
	return ctx.JSON(msg)
}	
