// Package routes contain url paths mapping to controllers
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/controllers"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/middlewares"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	v1 := router.Group("/api/v1")
	v1.Use(middlewares.AuthMiddleWare())

	v1.GET("/me", controllers.GetUser(client))
	v1.PUT("/me", controllers.UpdateUser(client))

	v1.GET("/movie/:imdb_id", controllers.GetMovie(client))
	v1.POST("/addmovie", controllers.AddMovie(client))
	v1.GET("/recommendedmovies", controllers.GetRecommendedMovies(client))
	v1.PATCH("/movie/:imdb_id/updatereview", controllers.AdminReviewUpdate(client))
}
