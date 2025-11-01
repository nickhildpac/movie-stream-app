package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/controllers"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupUnProtectedRoutes(router *gin.Engine, client *mongo.Client, mailChan chan models.MailData) {
	v1 := router.Group("/api/v1")
	v1.GET("/movies", controllers.GetMovies(client))
	v1.POST("/register", controllers.RegisterUser(client))
	v1.POST("/login", controllers.LoginUser(client))
	v1.POST("/logout", controllers.LogoutHandler(client))
	v1.GET("/genres", controllers.GetGenres(client))
	v1.POST("/refresh", controllers.RefreshTokenHandler(client))
	v1.POST("/request-reset", controllers.RequestResetPassword(client, mailChan))
	v1.POST("/reset-password", controllers.ResetPassword(client))
	v1.GET("/auth/google/login", controllers.GoogleLogin(client))
	v1.GET("/auth/google/callback", controllers.GoogleCallback(client))
}
