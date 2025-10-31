package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/database"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/models"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func RegisterUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		usersCollection := database.OpenCollection("users", client)

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input data"})
		}
		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid credentials", "details": err.Error()})
			return

		}
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()
		count, err := usersCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User already exists"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		user.UserID = bson.NewObjectID().Hex()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.Password = hashedPassword

		result, err := usersCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		c.JSON(http.StatusCreated, result)
	}
}

func LoginUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userLogin models.UserLogin

		if err := c.ShouldBindJSON(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalide input data"})
			return
		}

		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		userCollection := database.OpenCollection("users", client)

		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.D{{Key: "email", Value: userLogin.Email}}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userLogin.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		token, refreshToken, err := utils.GenerateAllTokens(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.Role, foundUser.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		err = utils.UpdateAllTokens(foundUser.UserID, token, refreshToken, client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tokens"})
			return
		}
		http.SetCookie(c.Writer, &http.Cookie{
			Name:  "access_token",
			Value: token,
			Path:  "/",
			// Domain:   "localhost",
			MaxAge:   86400,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})
		http.SetCookie(c.Writer, &http.Cookie{
			Name:  "refresh_token",
			Value: refreshToken,
			Path:  "/",
			// Domain:   "localhost",
			MaxAge:   604800,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		c.JSON(http.StatusOK, models.UserResponse{
			UserId:    foundUser.UserID,
			FirstName: foundUser.FirstName,
			LastName:  foundUser.LastName,
			Email:     foundUser.Email,
			Role:      foundUser.Role,
			Token:     token,
			// RefreshToken:    refreshToken,
			FavouriteGenres: foundUser.FavouriteGenres,
		})
	}
}

func LogoutHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Clear the access_token cookie

		var UserLogout struct {
			UserID string `json:"user_id"`
		}

		err := c.ShouldBindJSON(&UserLogout)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		fmt.Println("User ID from Logout request:", UserLogout.UserID)

		err = utils.UpdateAllTokens(UserLogout.UserID, "", "", client) // Clear tokens in the database
		// Optionally, you can also remove the user session from the database if needed
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error logging out"})
			return
		}
		// c.SetCookie(
		// 	"access_token",
		// 	"",
		// 	-1, // MaxAge negative → delete immediately
		// 	"/",
		// 	"localhost", // Adjust to your domain
		// 	true,        // Use true in production with HTTPS
		// 	true,        // HttpOnly
		// )
		http.SetCookie(c.Writer, &http.Cookie{
			Name:  "access_token",
			Value: "",
			Path:  "/",
			// Domain:   "localhost",
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		// // Clear the refresh_token cookie
		// c.SetCookie(
		// 	"refresh_token",
		// 	"",
		// 	-1,
		// 	"/",
		// 	"localhost",
		// 	true,
		// 	true,
		// )
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}

func RefreshTokenHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			fmt.Println("error", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve refresh token from cookie"})
			return
		}

		claim, err := utils.ValidateRefreshToken(refreshToken)
		if err != nil || claim == nil {
			fmt.Println("error", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			return
		}

		userCollection := database.OpenCollection("users", client)

		var user models.User
		err = userCollection.FindOne(ctx, bson.D{{Key: "user_id", Value: claim.UserID}}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		newToken, newRefreshToken, _ := utils.GenerateAllTokens(user.Email, user.FirstName, user.LastName, user.Role, user.UserID)
		err = utils.UpdateAllTokens(user.UserID, newToken, newRefreshToken, client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating tokens"})
			return
		}

		c.SetCookie("access_token", newToken, 86400, "/", "localhost", true, true)          // expires in 24 hours
		c.SetCookie("refresh_token", newRefreshToken, 604800, "/", "localhost", true, true) // expires in 1 week

		c.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed", "access_token": newToken})
	}
}

func RequestResetPassword(client *mongo.Client, mailChan chan models.MailData) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email" validate:"required,email"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		usersCollection := database.OpenCollection("users", client)
		var user models.User
		err := usersCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		token, err := utils.GeneratePasswordResetToken(user.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
			return
		}

		_, err = usersCollection.UpdateOne(
			ctx,
			bson.M{"user_id": user.UserID},
			bson.M{"$set": bson.M{
				"password_reset_token":   token,
				"password_reset_expires": time.Now().Add(time.Minute * 15),
			}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user with reset token"})
			return
		}

		// Assuming you have a function to send emails
		resetLink := fmt.Sprintf("http://localhost:5173/reset-password?token=%s", token)
		mailData := models.MailData{
			To:       user.Email,
			From:     "no-reply@movieapp.com",
			Subject:  "Password Reset",
			Content:  resetLink,
			Template: "password-reset.html",
		}
		// This is a placeholder for your email sending logic
		// You would typically call a utility function here, e.g., utils.SendEmail(mailData)
		fmt.Println("Sending password reset email:", mailData)
		mailChan <- mailData

		c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
	}
}

func ResetPassword(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Token       string `json:"token" validate:"required"`
			NewPassword string `json:"new_password" validate:"required,min=6"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		usersCollection := database.OpenCollection("users", client)
		var user models.User
		err := usersCollection.FindOne(ctx, bson.M{"password_reset_token": req.Token}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired token"})
			return
		}

		if time.Now().After(user.PasswordResetExpires) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token has expired"})
			return
		}

		hashedPassword, err := HashPassword(req.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		_, err = usersCollection.UpdateOne(
			ctx,
			bson.M{"user_id": user.UserID},
			bson.M{"$set": bson.M{
				"password":               hashedPassword,
				"password_reset_token":   "",
				"password_reset_expires": time.Time{},
			}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
	}
}
