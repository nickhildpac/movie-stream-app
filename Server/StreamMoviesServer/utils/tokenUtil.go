// Package utils consists function used for tokens and auth
package utils

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Role      string
	UserID    string
	jwt.RegisteredClaims
}

var (
	SecretKey        string = os.Getenv("SECRET_KEY")
	SecretRefreshKey string = os.Getenv("SECRET_REFRESH_KEY")
)

func GenerateAllTokens(email, firstName, lastName, role, userID string) (string, string, error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserID:    userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MovieStreamApp",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", "", err
	}

	refreshClaims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserID:    userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(SecretRefreshKey))
	if err != nil {
		return "", "", err
	}
	return signedToken, signedRefreshToken, nil
}

func GeneratePasswordResetToken(userID string) (string, error) {
	claims := &SignedDetails{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MovieStreamApp",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func UpdateAllTokens(userID, token, refreshToken string, client *mongo.Client) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateData := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"update_at":     updateAt,
		},
	}

	userCollection := database.OpenCollection("users", client)

	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userID}, updateData)
	if err != nil {
		return err
	}
	return nil
}

func GetAccessToken(c *gin.Context) (string, error) {
	// authHeader := c.Request.Header.Get("Authorization")
	// if authHeader == "" {
	// 	return "", errors.New("Authorization header is required")
	// }
	// tokenString := authHeader[len("Bearer "):]

	// if tokenString == "" {
	// 	return "", errors.New("Bearer token is required")
	// }
	tokenString, err := c.Cookie("access_token")
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}

func GetUserIDFromContext(c *gin.Context) (string, error) {
	userID, exists := c.Get("userID")

	if !exists {
		return "", errors.New("userId does not exists in this context")
	}

	id, ok := userID.(string)

	if !ok {
		return "", errors.New("unable to retrieve userId")
	}

	return id, nil
}

func GetRoleFromContext(c *gin.Context) (string, error) {
	role, exists := c.Get("role")
	if !exists {
		return "", errors.New("role does not exists in this context")
	}
	memberRole, ok := role.(string)
	if !ok {
		return "", errors.New("unable to retrieve userId")
	}

	return memberRole, nil
}

func ValidateRefreshToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretRefreshKey), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token has expired")
	}

	return claims, nil
}
