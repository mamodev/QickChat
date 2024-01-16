package api

import (
	"log"
	"server/db"

	"github.com/gofiber/fiber/v2"
)

func ChatGet(c *fiber.Ctx) error {

	c.Accepts("application/json")

	userId := c.Locals("userID").(string)


	conn, err := db.Conn(c.Context())
	if err != nil {
		log.Println("Error connecting to database: ", err)
		return fiber.ErrInternalServerError
	}
	defer conn.Close()

	rows, err := conn.QueryContext(c.Context(), `
	SELECT cu.chat_id, 
	CASE WHEN cc.count = 2 
		THEN (SELECT users.username FROM chat_user LEFT JOIN users ON chat_user.user_id = users.id WHERE chat_user.chat_id = c.id AND chat_user.user_id <> $1)
		ELSE c.name
	END AS name, 
	
	CASE WHEN cc.count = 2
		THEN (SELECT users.profile_picture FROM chat_user LEFT JOIN users ON chat_user.user_id = users.id WHERE chat_user.chat_id = c.id AND chat_user.user_id <> $1)
		ELSE c.picture
	END AS picture

FROM (SELECT chat_id, count(*) as count FROM chat_user GROUP BY chat_id) as cc
LEFT JOIN chat_user cu ON cc.chat_id = cu.chat_id
LEFT JOIN users u ON cu.user_id = u.id
LEFT JOIN chat c ON cu.chat_id = c.id
WHERE cu.user_id = $1`, userId)

	if err != nil {
		log.Println("Error querying database: ", err)
		return fiber.ErrInternalServerError
	}

	defer rows.Close()

	type Chat struct {
		Id string `json:"id"`
		Name string `json:"name"`
		Picture string `json:"picture"`
	}

	chats := make([]Chat, 0)

	for rows.Next() {
		var chat Chat
		err := rows.Scan(&chat.Id, &chat.Name, &chat.Picture)
		if err != nil {
			log.Println("Error scanning rows: ", err)
			return fiber.ErrInternalServerError
		}

		chat.Picture = "/api/images/" + chat.Picture

		chats = append(chats, chat)
	}

	return c.JSON(chats)
}



// SELECT cu.chat_id, 
// 	CASE WHEN cc.count = 2 
// 		THEN (SELECT users.username FROM chat_user LEFT JOIN users ON chat_user.user_id = users.id WHERE chat_user.chat_id = c.id AND chat_user.user_id <> $1)
// 		ELSE c.name
// 	END AS name, 
	
// 	CASE WHEN cc.count = 2
// 		THEN (SELECT users.profile_picture FROM chat_user LEFT JOIN users ON chat_user.user_id = users.id WHERE chat_user.chat_id = c.id AND chat_user.user_id <> $1)
// 		ELSE c.picture
// 	END AS picture

// FROM (SELECT chat_id, count(*) as count FROM chat_user GROUP BY chat_id) as cc
// LEFT JOIN chat_user cu ON cc.chat_id = cu.chat_id
// LEFT JOIN users u ON cu.user_id = u.id
// LEFT JOIN chats c ON cu.chat_id = c.id
// WHERE cu.user_id = $1