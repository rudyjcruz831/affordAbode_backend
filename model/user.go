package model

import "time"

// User ...
type Users struct {
	ID        string     `json:"user_id" gorm:"column:user_id;primary_key"`
	Email     string     `json:"email" binding:"required,email" gorm:"column:email;unique;not null"`
	Username  string     `json:"username" binding:"required" gorm:"column:username;unique;not null"`
	Password  string     `json:"-" binding:"required,gte=6,lte=30" gorm:"column:password"`
	FirstName string     `json:"first_name" gorm:"column:first_name"`
	LastName  string     `json:"last_name" gorm:"column:last_name"`
	Role      string     `json:"-" gorm:"column:user_role"`
	CreatedOn time.Time  `json:"-" gorm:"column:created_on"`
	UpdatedAt time.Time  `json:"-" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"column:deleted_at"`
}

// Token Response ...
type Response struct {
	AccessToken string `json:"access_token"`
	ID          string `json:"user_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	TokenType   string `json:"token_type"`
}

// TokenID ...
type TokenID struct {
	ID        string `json:"sub"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	Email     string `json:"email"`
}
