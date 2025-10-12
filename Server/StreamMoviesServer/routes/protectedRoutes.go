package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/controllers"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/middlewares"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	router.Use(middlewares.AuthMiddleWare())

	router.GET("/movie/:imdb_id", controllers.GetMovie(client))
	router.POST("/addmovie", controllers.AddMovie(client))
	// router.GET("/recommendedmovies", controllers.GetRecommendedMovies(client))
	// router.PATCH("/updatereview/:imdb_id", controllers.AdminReviewUpdate(client))
}
