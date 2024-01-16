package api

import (
	"log"
	"server/db"

	"github.com/gofiber/fiber/v2"
)

var getFriendRequestsQuery = `
select fr.id, fr.sender_id,
	u.username as sender_name, u.email as sender_email, u.profile_picture as sender_profile_picture
from friend_request fr join users u on fr.sender_id = u.id
where fr.receiver_id = $1
order by fr.created_at desc, u.username asc
limit 40 offset $2 * 40
`

func FriendRequestGet (c *fiber.Ctx) error {
	
	c.Accepts("application/json")

	currentUser := c.Locals("userID").(string)
	page := c.Query("page")
	if page == "" {
		page = "0"
	}

	type friendRequest struct {
		Id string `json:"id"`
		SenderId string `json:"sender_id"`
		SenderName string `json:"sender_name"`
		SenderEmail string `json:"sender_email"`
		SenderProfilePicture string `json:"sender_profile_picture"`
	}
		
	var friendRequests []friendRequest = make([]friendRequest, 0)

	conn, err := db.Conn(c.Context())
	if err != nil {
		log.Println("Error connecting to database:", err)
		return fiber.ErrInternalServerError
	}
	defer conn.Close()

	rows, err := conn.QueryContext(c.Context(), getFriendRequestsQuery, currentUser, page)

	if err != nil {
		log.Println("Error querying database:", err)
		return fiber.ErrInternalServerError
	}
	defer rows.Close()

	for rows.Next() {
		var friendRequest friendRequest

		err := rows.Scan(&friendRequest.Id, &friendRequest.SenderId, &friendRequest.SenderName, &friendRequest.SenderEmail, &friendRequest.SenderProfilePicture)
		if err != nil {
			log.Println("Error scanning row:", err)
			return fiber.ErrInternalServerError
		}

		friendRequest.SenderProfilePicture = "/api/images/" + friendRequest.SenderProfilePicture
		friendRequests = append(friendRequests, friendRequest)
	}

	return c.JSON(friendRequests)
}