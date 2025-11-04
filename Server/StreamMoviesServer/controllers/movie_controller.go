// Package controllers contains movies and users helper functions
package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/database"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/models"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/utils"
	"github.com/tmc/langchaingo/llms/openai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func GetGenres(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()
		// var movieCollection *mongo.Collection = database.OpenCollection("movies")
		genresCollection := database.OpenCollection("genres", client)

		cursor, err := genresCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch genres"})
			return
		}
		defer cursor.Close(ctx)

		var results []models.Genre
		if err = cursor.All(ctx, &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode genres"})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}

func GetMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()
		movieCollection := database.OpenCollection("movies", client)
		var movies []models.Movie
		cursor, err := movieCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)
		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies"})
			return
		}
		c.JSON(http.StatusOK, movies)
	}
}

func GetMovie(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()
		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
			return
		}
		var movie models.Movie
		movieCollection := database.OpenCollection("movies", client)
		err := movieCollection.FindOne(ctx, bson.D{{Key: "imdb_id", Value: movieID}}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get movies"})
			return
		}
		c.JSON(http.StatusOK, movie)
	}
}

func AddOrUpdateGenre(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()
		var genre models.Genre
		genresCollection := database.OpenCollection("genres", client)
		if err := c.ShouldBindJSON(&genre); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		filter := bson.M{"genre_id": genre.GenreID}
		update := bson.M{
			"$set": bson.M{
				"genre_name": genre.GenreName,
			},
		}
		res := genresCollection.FindOneAndUpdate(ctx, filter, update)
		if res.Err() == mongo.ErrNoDocuments {
			_, err := genresCollection.InsertOne(ctx, genre)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusAccepted, gin.H{"message": "genre created"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"message": "genre updated"})
	}
}
func AddMovie(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()
		var movie models.Movie
		movieCollection := database.OpenCollection("movies", client)
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res := movieCollection.FindOne(ctx, bson.M{"imdb_id": movie.ImdbID})
		if res.Err() == mongo.ErrNoDocuments {
			_, err := movieCollection.InsertOne(ctx, movie)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

		}
		var insertedMovie models.Movie
		err := movieCollection.FindOne(ctx, bson.D{{Key: "imdb_id", Value: movie.ImdbID}}).Decode(&insertedMovie)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"state": "success", "message": "posted data", "data": insertedMovie})
	}
}

func AdminReviewUpdate(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, err := utils.GetRoleFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role not found in context"})
			return
		}
		if role != "ADMIN" && role != "USER" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be admin or user"})
			return
		}
		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie Id required"})
			return
		}
		var req struct {
			AdminReview string `json:"admin_review"`
		}
		var resp struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		sentiment, rankVal, err := GetReviewRanking(req.AdminReview, client, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting review ranking"})
			return
		}
		filter := bson.M{"imdb_id": movieID}
		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_value": rankVal,
					"ranking_name":  sentiment,
				},
			},
		}
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()
		movieCollection := database.OpenCollection("movies", client)
		result, err := movieCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating movie"})
			return
		}
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
		resp.AdminReview = req.AdminReview
		resp.RankingName = sentiment
		c.JSON(http.StatusOK, resp)
	}
}

func GetReviewRanking(adminReview string, client *mongo.Client, c *gin.Context) (string, int, error) {
	rankings, err := GetRankings(client, c)
	if err != nil {
		return "", 0, err
	}
	sentimentDelimited := ""
	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimited = sentimentDelimited + ranking.RankingName + ","
		}
		sentimentDelimited = strings.Trim(sentimentDelimited, ",")
	}
	err = godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}
	OpenAiAPIKey := os.Getenv("OPENAI_API_KEY")
	if OpenAiAPIKey == "" {
		return "", 0, errors.New("could not get open ai api key")
	}
	llm, err := openai.New(openai.WithToken(OpenAiAPIKey))
	if err != nil {
		return "", 0, err
	}
	basePromptTemplate := os.Getenv("BASE_PROMPT_TEMPLATE")
	basePrompt := strings.Replace(basePromptTemplate, "{rankings}", sentimentDelimited, 1)
	response, err := llm.Call(c, basePrompt+adminReview)
	if err != nil {
		return "", 0, nil
	}
	rankVal := 0
	for _, rankings := range rankings {
		if rankings.RankingName == response {
			rankVal = rankings.RankingValue
			break
		}
	}
	return response, rankVal, nil
}

func GetRankings(client *mongo.Client, c *gin.Context) ([]models.Ranking, error) {
	var rankings []models.Ranking
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()
	rankingCollection := database.OpenCollection("rankings", client)
	cursor, err := rankingCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}
	return rankings, nil
}

func GetRecommendedMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		favouriteGenres, err := GetUsersFavouriteGenres(userID, client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = godotenv.Load(".env")
		if err != nil {
			log.Println("Warning: .env file not found")
		}
		var recommendedMoviesLimitVal int64 = 5
		recommendedMoviesLimitValStr := os.Getenv("RECOMMENDED_MOVIE_LIMIT")
		if recommendedMoviesLimitValStr != "" {
			recommendedMoviesLimitVal, _ = strconv.ParseInt(recommendedMoviesLimitValStr, 10, 64)
		}
		findOptions := options.Find()
		findOptions.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})
		findOptions.SetLimit(recommendedMoviesLimitVal)
		filter := bson.M{"genre.genre_name": bson.M{"$in": favouriteGenres}}
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()
		movieCollection := database.OpenCollection("movies", client)
		cursor, err := movieCollection.Find(ctx, filter, findOptions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recommended movies"})
			return
		}
		defer cursor.Close(ctx)
		var recommendedMovies []models.Movie
		if err := cursor.All(ctx, &recommendedMovies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, recommendedMovies)
	}
}

func GetUsersFavouriteGenres(userID string, client *mongo.Client) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	filter := bson.M{"user_id": userID}
	projection := bson.M{
		"favourite_genres.genre_name": 1,
		"_id":                         0,
	}
	opts := options.FindOne().SetProjection(projection)
	var result bson.M
	usersCollection := database.OpenCollection("users", client)
	err := usersCollection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		}
	}
	favGenresArray, ok := result["favourite_genres"].(bson.A)
	if !ok {
		return []string{}, errors.New("unable to retrieve favourite_genres for user")
	}
	var genreNames []string
	for _, item := range favGenresArray {
		if genreMap, ok := item.(bson.D); ok {
			for _, elem := range genreMap {
				if elem.Key == "genre_name" {
					if name, ok := elem.Value.(string); ok {
						genreNames = append(genreNames, name)
					}
				}
			}
		}
	}
	return genreNames, nil
}
