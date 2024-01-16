package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"server/db"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SubscriptionInfo struct {
	Endpoint string `json:"endpoint" validate:"required"`
}	

func SubscribeHandler(c *fiber.Ctx) error {

	c.Accepts("application/json")

	var subscriptionInfo SubscriptionInfo

	userID := c.Locals("userID").(string)

	if err := c.BodyParser(&subscriptionInfo); err != nil {
		return c.Status(http.StatusBadRequest).SendString("Failed to decode JSON")
	}

	if err:= Validator.Struct(subscriptionInfo); err != nil {
		return c.Status(http.StatusBadRequest).SendString("Failed to validate JSON")
	}

	conn, err := db.Conn(c.Context())
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	defer conn.Close()


	id := uuid.New().String()
	_, err = conn.ExecContext(c.Context(), "INSERT INTO push_subscription(id, user_id, endpoint) VALUES(?, ?, ?)", id, userID, subscriptionInfo.Endpoint)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"id": id,
	})
}

func UnsubscribeHandler(c *fiber.Ctx) error {
	
	c.Accepts("application/json")

	id := c.Params("id")

	conn, err := db.Conn(c.Context())
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	defer conn.Close()

	_, err = conn.ExecContext(c.Context(), "DELETE FROM push_subscription WHERE id = ?", id)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func SendNotification (c *fiber.Ctx, msg, user string) error {

	conn, err := db.Conn(c.Context())
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}	

	fmt.Println("Sending notification to", user)

	defer conn.Close()

	rows, err := conn.QueryContext(c.Context(), "SELECT endpoint FROM push_subscription WHERE user_id = ?", user)

	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	defer rows.Close()

	endpoints := make([]string, 0)
	for rows.Next() {
		var endpoint string
		err = rows.Scan(&endpoint)
		if err != nil {
			log.Println(err)
			return fiber.ErrInternalServerError
		}

		endpoints = append(endpoints, endpoint)
	}

	for _, endpoint := range endpoints {

		// vapidDetails := webpush.VapidDetails{
		// 	Subject: "mailto: <marco.morozzi2002@gmail.com>",
		// 	PublicKey: "BKrFuHA5xh5D-67Mpl6_t81pR_ezyNMWf3AYLfagOZKUK6kjoyjc1dpNgWkyfqh7e-q2oety0B68XxUEcLWAirw",
		// 	PrivateKey: "X908hc9vb8gQOu-6AokEJD8xR04-gXX8H-LvExQs6VA",
		// }

		log.Println("Sending notification to", endpoint)

		body := "{ \"title\": \"Chat App\", \"body\": \"" + msg + "\" }"


		req, err := http.NewRequest("POST", endpoint, strings.NewReader(body))

		if err != nil {
			log.Println(err)
			return fiber.ErrInternalServerError
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "key=X908hc9vb8gQOu-6AokEJD8xR04-gXX8H-LvExQs6VA")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			return fiber.ErrInternalServerError
		}

		defer resp.Body.Close()


		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return fiber.ErrInternalServerError
		}

		log.Println("Sent notification to", resp.StatusCode, string(respBody))
	}

	return nil
}