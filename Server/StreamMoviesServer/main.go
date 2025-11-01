package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/controllers"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/database"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/models"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/routes"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/utils"
)

func main() {
	router := gin.Default()

	router.GET("/healthcheck", func(c *gin.Context) {
		c.String(200, "ok")
	})

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to find .env file")
	}

	controllers.InitGoogleOAuth()

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	var origins []string
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
			log.Println("Allowed Origin:", origins[i])
		}
	} else {
		origins = []string{"http://localhost:5173"}
		log.Println("Allowed Origin: http://localhost:5173")
	}

	config := cors.Config{}
	config.AllowOrigins = origins
	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	// config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))
	router.Use(gin.Logger())

	client := database.Connect()

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to reach server: %v", err)
	}
	defer func() {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()
	mailChan := make(chan models.MailData)
	defer close(mailChan)
	utils.ListenForMail(mailChan)
	// msg := models.MailData{
	// 	To:       "oH4yH@example.com",
	// 	From:     "me@example.com",
	// 	Subject:  "Test Email",
	// 	Content:  "Test Reset link in Email",
	// 	Template: "email-template.html",
	// }
	// mailChan <- msg

	routes.SetupUnProtectedRoutes(router, client, mailChan)
	routes.SetupProtectedRoutes(router, client)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
