package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/controllers"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupUnProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	v1 := router.Group("/v1")
	v1.GET("/movies", controllers.GetMovies(client))
	v1.POST("/register", controllers.RegisterUser(client))
	v1.POST("/login", controllers.LoginUser(client))
	v1.POST("/logout", controllers.LogoutHandler(client))
	// router.GET("/genres", controllers.GetGenres(client))
	v1.POST("/refresh", controllers.RefreshTokenHandler(client))
}
