package api

import (
	"log"
	"server/db"

	"github.com/gofiber/fiber/v2"
)


var getUserQuery = `
select 
	u.id,
	username,
	email,
	profile_picture,
	case when fr.id is null then false else true end as friend_request_sent,
	(select count(*) from friend f where f.user1_id = $1 and f.user2_id = u.id) > 0 as is_friend,
	fr.id as friend_request
from users u left join friend_request fr on fr.sender_id = $1 and fr.receiver_id = u.id
where (lower(username) like $2 or lower(email) like $2) and u.id != $1
order by is_friend desc, friend_request_sent desc, username
limit 40 offset $3 * 40
`

func Users (c *fiber.Ctx) error {

	c.Accepts("application/json")

	// get user from url param
	user := c.Query("user")
	currentUser := c.Locals("userID").(string)

	page := c.Query("page")
	if page == "" {
		page = "0"
	}

	// log all raw url with query params
	conn, err := db.Conn(c.Context())
	if err != nil {
		log.Println("Error connecting to database:", err)
		return fiber.ErrInternalServerError
	}

	defer conn.Close()

	rows, err := conn.QueryContext(c.Context(), getUserQuery, currentUser, "%" + user + "%", page)
	defer rows.Close()
	if err != nil {
		log.Println("Error querying database:", err)
		return fiber.ErrInternalServerError
	}

	type User struct {
		Id string `json:"id"`
		Username string `json:"username"`
		Email string `json:"email"`
		ProfilePicture string `json:"profile_picture"`
		FriendRequestSent bool `json:"friend_request_sent"`
		FriendRequestId *string `json:"friend_request_id"`
		IsFriend bool `json:"is_friend"`
	}

	users := make([]User, 0)

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.ProfilePicture, &user.FriendRequestSent, &user.IsFriend, &user.FriendRequestId)
		if err != nil {
			log.Println("Error scanning rows: ", err)
			return fiber.ErrInternalServerError
		}

		user.ProfilePicture = "/api/images/" + user.ProfilePicture

		users = append(users, user)
	}

	return c.JSON(users)
}