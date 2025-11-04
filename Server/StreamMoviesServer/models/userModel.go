// Package models includes domain used in the application
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID                   bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID               string        `json:"user_id" bson:"user_id"`
	FirstName            string        `json:"first_name" bson:"first_name" validate:"required,min=2,max=100"`
	LastName             string        `json:"last_name" bson:"last_name" validate:"required,min=2,max=100"`
	Email                string        `json:"email" bson:"email" validate:"required,email"`
	Password             string        `json:"password,omitempty" bson:"password,omitempty" validate:"omitempty,min=6"`
	Role                 string        `json:"role" bson:"role" validate:"oneof=ADMIN USER"`
	CreatedAt            time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt            time.Time     `json:"update_at" bson:"update_at"`
	Token                string        `json:"token" bson:"token"`
	RefreshToken         string        `json:"refresh_token" bson:"refresh_token"`
	FavouriteGenres      []Genre       `json:"favourite_genres" bson:"favourite_genres" validate:"required,dive"`
	PasswordResetToken   string        `json:"password_reset_token,omitempty" bson:"password_reset_token,omitempty"`
	PasswordResetExpires time.Time     `json:"password_reset_expires,omitzero" bson:"password_reset_expires,omitzero"`
	AuthProvider         string        `json:"auth_provider" bson:"auth_provider"`
}
type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
type UserResponse struct {
	UserID          string  `json:"user_id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Email           string  `json:"email"`
	Role            string  `json:"role"`
	Token           string  `json:"token"`
	RefreshToken    string  `json:"refresh_token"`
	FavouriteGenres []Genre `json:"favourite_genres"`
}

type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}

type UpdateUser struct {
	UserID          string    `json:"user_id" bson:"user_id"`
	FirstName       string    `json:"first_name" bson:"first_name" validate:"required,min=2,max=100"`
	LastName        string    `json:"last_name" bson:"last_name" validate:"required,min=2,max=100"`
	Email           string    `json:"email" bson:"email" validate:"required,email"`
	UpdatedAt       time.Time `json:"update_at" bson:"update_at"`
	FavouriteGenres []Genre   `json:"favourite_genres" bson:"favourite_genres" validate:"required,dive"`
}

type LogoutRequest struct {
	UserID string `json:"user_id"`
}

type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type PasswordReset struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}
