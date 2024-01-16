package api

import (
	"database/sql"
	"log"
	"server/db"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ChatCreateRequest struct {
	Users []string `json:"users" validate:"required"`
	Name string `json:"name" validate:"required"`
}
	
func ChatCreate(c *fiber.Ctx) error {

	c.Accepts("application/json")

	userId := c.Locals("userID").(string)

	var body ChatCreateRequest
	err := c.BodyParser(&body)
	if err != nil {
		log.Println("Error parsing body: ", err)
		return fiber.ErrBadRequest
	}

	if len(body.Users) == 0 {
		return fiber.ErrBadRequest
	}

	conn, err := db.Conn(c.Context())
	if err != nil {
		log.Println("Error connecting to database: ", err)
		return fiber.ErrInternalServerError
	}
	defer conn.Close()

	tx, err := conn.BeginTx(c.Context(), nil)
	if err != nil {
		log.Println("Error beginning transaction: ", err)
		return fiber.ErrInternalServerError
	}

	shouldRollback := true
	defer func() {
		if shouldRollback {
			tx.Rollback()
		}
	}()

	chatId := uuid.New().String()

	// check all users exist
	userSqlStringArray := "("
	for i, user := range body.Users {
		if i != 0 {
			userSqlStringArray += ", "
		}
		userSqlStringArray += "'" + user + "'"
	}
	userSqlStringArray += ")"

	rows, err := tx.QueryContext(c.Context(), "select id from users join friend on users.id = friend.user2_id where user1_id = $1 and id in " + userSqlStringArray, userId)
	
	if err != nil {
		log.Println("Error querying database: ", err)
		return fiber.ErrInternalServerError
	}
	defer rows.Close()

	userIds := make([]string, 0)
	for rows.Next() {
		var userId string
		err := rows.Scan(&userId)
		if err != nil {
			log.Println("Error scanning rows: ", err)
			return fiber.ErrInternalServerError
		}

		userIds = append(userIds, userId)
	}

	if len(userIds) != len(body.Users) {
		return fiber.ErrBadRequest
	}

	// insert chat
	_, err = tx.ExecContext(c.Context(), "insert into chat (id, name) values ($1, $2)", chatId, body.Name)
	if err != nil {
		log.Println("Error inserting chat: ", err)
		return fiber.ErrInternalServerError
	}

	userIds = append(userIds, userId)

	if len(userIds) == 2 {
		// check if there is already a chat between the two users
		var existingChatID string
		err = tx.QueryRowContext(c.Context(), `SELECT chat_id
			FROM chat_user
			WHERE user_id = $1 AND chat_id IN (
				SELECT chat_id
				FROM chat_user
				WHERE user_id = $2
			)
			LIMIT 1
		`, userIds[0], userIds[1]).Scan(&existingChatID)
	
		if err != nil && err != sql.ErrNoRows  {
			log.Println("Error querying database: ", err)
			return fiber.ErrInternalServerError
		}
		
		if err == nil {
			// SendNotification(c, "Someone added you to a chat", userIds[1])
			return c.JSON(fiber.Map{
				"id": existingChatID,
			})
		}
	}
	
	// insert chat_user
	for _, userId := range userIds {
		_, err = tx.ExecContext(c.Context(), "insert into chat_user (chat_id, user_id) values ($1, $2)", chatId, userId)
		if err != nil {
			log.Println("Error inserting chat_user: ", err)
			return fiber.ErrInternalServerError
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction: ", err)
		return fiber.ErrInternalServerError
	}

	shouldRollback = false

	return c.JSON(fiber.Map{
		"id": chatId,
	})
}