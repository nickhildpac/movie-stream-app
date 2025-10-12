package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/controllers"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/middlewares"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	v1 := router.Group("/v1")
	v1.Use(middlewares.AuthMiddleWare())

	v1.GET("/movie/:imdb_id", controllers.GetMovie(client))
	v1.POST("/addmovie", controllers.AddMovie(client))
	// router.GET("/recommendedmovies", controllers.GetRecommendedMovies(client))
	// router.PATCH("/updatereview/:imdb_id", controllers.AdminReviewUpdate(client))
}
