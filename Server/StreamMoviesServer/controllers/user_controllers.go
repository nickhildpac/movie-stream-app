package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/database"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/models"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config

func InitGoogleOAuth() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func GoogleLogin(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var state string
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		state = base64.URLEncoding.EncodeToString(b)

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "oauthstate",
			Value:    state,
			Expires:  time.Now().Add(20 * time.Minute),
			HttpOnly: true,
		})
		url := googleOauthConfig.AuthCodeURL(state)
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func GoogleCallback(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		oauthState, _ := c.Cookie("oauthstate")
		if c.Query("state") != oauthState {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
			return
		}
		code := c.Query("code")
		token, err := googleOauthConfig.Exchange(context.Background(), code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token", "details": err.Error()})
			return
		}

		response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info", "details": err.Error()})
			return
		}
		defer response.Body.Close()

		contents, err := io.ReadAll(response.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read user info", "details": err.Error()})
			return
		}
		var userInfo struct {
			Email     string `json:"email"`
			FirstName string `json:"given_name"`
			LastName  string `json:"family_name"`
		}
		if err := json.Unmarshal(contents, &userInfo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info", "details": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		usersCollection := database.OpenCollection("users", client)
		var user models.User
		err = usersCollection.FindOne(ctx, bson.M{"email": userInfo.Email}).Decode(&user)

		if err == mongo.ErrNoDocuments {
			// User does not exist, create a new one
			user = models.User{
				UserID:       bson.NewObjectID().Hex(),
				FirstName:    userInfo.FirstName,
				LastName:     userInfo.LastName,
				Email:        userInfo.Email,
				Role:         "USER",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				AuthProvider: "google",
			}

			_, err := usersCollection.InsertOne(ctx, user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
			return
		}
		// User exists or was just created, generate tokens
		appToken, refreshToken, err := utils.GenerateAllTokens(user.Email, user.FirstName, user.LastName, user.Role, user.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		err = utils.UpdateAllTokens(user.UserID, appToken, refreshToken, client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tokens"})
			return
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "access_token",
			Value:    appToken,
			Path:     "/",
			MaxAge:   86400,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Path:     "/",
			MaxAge:   604800,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:5173") // Redirect to frontend
	}
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.User true "User object"
// @Success 201 {object} primitive.ObjectID
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /register [post]
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
		user.AuthProvider = "local"

		result, err := usersCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		c.JSON(http.StatusCreated, result)
	}
}

// LoginUser godoc
// @Summary Login a user
// @Description Login a user with email and password
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.UserLogin true "User login object"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /login [post]
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

		if foundUser.AuthProvider != "local" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Please sign in with " + foundUser.AuthProvider})
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
			Name:     "access_token",
			Value:    token,
			Path:     "/",
			MaxAge:   86400,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Path:     "/",
			MaxAge:   604800,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		c.JSON(http.StatusOK, models.UserResponse{
			UserID:          foundUser.UserID,
			FirstName:       foundUser.FirstName,
			LastName:        foundUser.LastName,
			Email:           foundUser.Email,
			Role:            foundUser.Role,
			Token:           token,
			FavouriteGenres: foundUser.FavouriteGenres,
		})
	}
}

// GetUser godoc
// @Summary Get user details
// @Description Get details of the currently logged in user
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /me [get]
func GetUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}

		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		userCollection := database.OpenCollection("users", client)

		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&foundUser)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, models.UserResponse{
			UserID:          foundUser.UserID,
			FirstName:       foundUser.FirstName,
			LastName:        foundUser.LastName,
			Email:           foundUser.Email,
			Role:            foundUser.Role,
			Token:           foundUser.Token,
			RefreshToken:    foundUser.RefreshToken,
			FavouriteGenres: foundUser.FavouriteGenres,
		})
	}
}

// LogoutHandler godoc
// @Summary Logout a user
// @Description Logout the current user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.LogoutRequest true "User ID"
// @Success 200 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /logout [post]
func LogoutHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userLogout models.LogoutRequest

		err := c.ShouldBindJSON(&userLogout)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		fmt.Println("User ID from Logout request:", userLogout.UserID)

		err = utils.UpdateAllTokens(userLogout.UserID, "", "", client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error logging out"})
			return
		}
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})
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

// RefreshTokenHandler godoc
// @Summary Refresh access token
// @Description Refresh the access token using the refresh token
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /refresh [post]
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

		c.SetCookie("access_token", newToken, 86400, "/", "localhost", true, true)
		c.SetCookie("refresh_token", newRefreshToken, 604800, "/", "localhost", true, true)

		c.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed", "access_token": newToken})
	}
}

// RequestResetPassword godoc
// @Summary Request a password reset
// @Description Request a password reset email
// @Tags users
// @Accept  json
// @Produce  json
// @Param email body models.PasswordResetRequest true "User email"
// @Success 200 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /request-reset [post]
func RequestResetPassword(client *mongo.Client, mailChan chan models.MailData) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.PasswordResetRequest
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

		resetLink := fmt.Sprintf("http://localhost:5173/reset-password?token=%s", token)
		mailData := models.MailData{
			To:       user.Email,
			From:     "no-reply@movieapp.com",
			Subject:  "Password Reset",
			Content:  resetLink,
			Template: "password-reset.html",
		}

		mailChan <- mailData
		fmt.Println("Sending password reset email:", mailData)

		c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
	}
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset the user's password
// @Tags users
// @Accept  json
// @Produce  json
// @Param password body models.PasswordReset true "Reset password object"
// @Success 200 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /reset-password [post]
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

// UpdateUser godoc
// @Summary Update user details
// @Description Update details of the currently logged in user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.UpdateUser true "User object"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /me [put]
func UpdateUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		log.Println("this is user Id ", userID)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}

		var updateData models.UpdateUser
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		usersCollection := database.OpenCollection("users", client)

		updateFields := bson.M{
			"updated_at":       time.Now(),
			"first_name":       updateData.FirstName,
			"email":            updateData.Email,
			"last_name":        updateData.LastName,
			"favourite_genres": updateData.FavouriteGenres,
		}

		result, err := usersCollection.UpdateOne(
			ctx,
			bson.M{"user_id": userID},
			bson.M{"$set": updateFields},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "details": err.Error()})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		var updatedUser models.User
		err = usersCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&updatedUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user"})
			return
		}

		c.JSON(http.StatusOK, models.UserResponse{
			UserID:          updatedUser.UserID,
			FirstName:       updatedUser.FirstName,
			LastName:        updatedUser.LastName,
			Email:           updatedUser.Email,
			Role:            updatedUser.Role,
			Token:           updatedUser.Token,
			RefreshToken:    updatedUser.RefreshToken,
			FavouriteGenres: updatedUser.FavouriteGenres,
		})
	}
}
