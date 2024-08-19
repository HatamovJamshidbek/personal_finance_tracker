package models

import "time"

type User struct {
	Id           string    `json:"id"`
	UserName     string    `json:"user_name"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Phone        string    `json:"phone"`
	Image        string    `json:"image"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
type Register struct {
	Username     string `json:"user_name"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Phone        string `json:"phone"`
	Image        string `json:"image"`
	Role         string `json:"role"`
}
type RegisterResponse struct {
	Id           string    `json:"id"`
	UserName     string    `json:"user_name"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Phone        string    `json:"phone"`
	Image        string    `json:"image"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}
type LoginRequest struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}
type LoginResponse struct {
	Id           string `json:"id"`
	UserName     string `json:"user_name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
}
type Token struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	ExpiredTime  float64 `json:"expired_time"`
}
type SendEmail struct {
	Email              string `json:"email"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

type SaveToken struct {
	UserId  string
	Token   string
	Revoked bool
}
