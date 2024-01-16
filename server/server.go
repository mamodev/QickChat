package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"server/api"
	"server/auth"
	"server/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)




func main() {
	err := db.Init("../database/database.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: api.ErrorHandler,
	})

	corsConfig := cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE",
	}


	reactApp := "../public"
	app.Static("/", reactApp)
	
	app.Use(cors.New(corsConfig))

	apiRoutes := app.Group("/api")
	apiRoutes.Post("/login", api.Login)
	apiRoutes.Post("/register", api.Register)
	
	authRoutes := apiRoutes.Group("");
	authRoutes.Use(auth.AuthMiddleware)

	authRoutes.Get("/chat", api.ChatGet);
	authRoutes.Post("/chat", api.ChatCreate);

	authRoutes.Get("/chat/:id/message", api.ChatMessageGet);
	authRoutes.Post("/chat/:id/message", api.ChatMessageCreate);
	authRoutes.Get("/chat/:id/avatar", api.ChatAvatarGet);

	authRoutes.Get("/user", api.Users);
	authRoutes.Get("/user/:id", api.UserGet);
	authRoutes.Get("/user/:id/avatar", api.UserAvatarGet);
	
	authRoutes.Post("/friend-request", api.FriendRequestCreate)
	authRoutes.Get("/friend-request", api.FriendRequestGet)
	authRoutes.Delete("/friend-request/:id", api.FriendRequestDelete)
	authRoutes.Post("/friend-request/:id/respond", api.FriendRequestRespond)
	
	authRoutes.Post("/logout", api.Logout)
	authRoutes.Get("/chats-update", api.ChatUpdate)

	authRoutes.Post("/push-subscription", api.SubscribeHandler)
	authRoutes.Delete("/push-subscription/:id", api.UnsubscribeHandler)

	authRoutes.Static("/images", "../database/img")


	go func () {
		app.ListenTLS(":3000", "../localhost.pem", "../localhost-key.pem")
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	fmt.Println("Server is shutting down...")
	app.ShutdownWithTimeout(2 * time.Second)
	// db.Close()
}
