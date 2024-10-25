package models

import (
	"time"
)

type UserRole string
type UserGender string

const (
	RoleUser     UserRole   = "user"
	RoleAdmin    UserRole   = "admin"
	GenderFemale UserGender = "female"
	GenderMale   UserGender = "male"
)

type UserQueryParams struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	FullName    string `json:"fullName"`
	DateOfBirth string `json:"dateOfBirth"`
	Role        string `json:"role"`
	Page        int    `json:"page" binding:"omitempty,min=1"`
	Limit       int    `json:"limit" binding:"omitempty,min=1,max=100"`
}

type CreateUserRequest struct {
	Email       string     `json:"email" validate:"required,email"`
	Username    string     `json:"username" validate:"required"`
	FullName    string     `json:"fullName" validate:"required"`
	Password    string     `json:"password" validate:"required,min=6"`
	Gender      UserGender `json:"gender" validate:"required"`
	Avatar      string     `json:"avatar" validate:"required"`
	DateOfBirth string     `json:"dateOfBirth" validate:"required"`
	Role        UserRole   `json:"role"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Message     string     `json:"message" validate:"required"`
	User        UserDetail `json:"user" validate:"required"`
	AccessToken string     `json:"accessToken" validate:"required"`
	ExpiresIn   int64      `json:"expiresIn" validate:"required"`
}

type UserDetail struct {
	ID          int        `json:"id" validate:"required"`
	Email       string     `json:"email" validate:"required"`
	Username    string     `json:"username" validate:"required"`
	FullName    string     `json:"fullName" validate:"required"`
	Gender      UserGender `json:"gender" validate:"required"`
	Avatar      string     `json:"avatar" validate:"required"`
	DateOfBirth string     `json:"dateOfBirth" validate:"required"`
	Role        UserRole   `json:"role" validate:"required"`
	CreatedAt   time.Time  `json:"createdAt" validate:"required"`
}

type UserToken struct {
	ID   int      `json:"id" validate:"required"`
	Role UserRole `json:"role" validate:"required"`
}

type UpdateUserResponse struct {
	Message string     `json:"message" validate:"required"`
	User    UserDetail `json:"user" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Password string `json:"password" validate:"required,min=6"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"accessToken" validate:"required"`
	ExpiresIn   int64  `json:"expiresIn" validate:"required"`
}

type PasswordUpdateRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=6"`
	NewPassword     string `json:"newPassword" validate:"required,min=6"`
}

type UserResponse struct {
	Data   []UserDetail `json:"data" validate:"required"`
	Paging PagingInfo   `json:"paging" validate:"required"`
}

type PagingInfo struct {
	Page       int `json:"page" validate:"required"`
	Limit      int `json:"limit" validate:"required"`
	TotalCount int `json:"totalCount" validate:"required"`
}
