package models

type CreateUserRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Username    string `json:"username" validate:"required"`
	FullName    string `json:"fullName" validate:"required"`
	Password    string `json:"password" validate:"required,min=6"`
	Gender      string `json:"gender" validate:"required"`
	Avatar      string `json:"avatar"`
	DateOfBirth string `json:"dateOfBirth"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Message     string     `json:"message" validate:"required"`
	User        UserDetail `json:"user" validate:"required"`
	AccessToken string     `json:"accessToken" validate:"required"`
}

type UserDetail struct {
	ID          int    `json:"id" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Username    string `json:"username" validate:"required"`
	FullName    string `json:"fullName" validate:"required"`
	Gender      string `json:"gender" validate:"required"`
	Avatar      string `json:"avatar"`
	DateOfBirth string `json:"dateOfBirth"`
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
	AccessToken string `json:"accessToken"`
}

type PasswordUpdateRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=6"`
	NewPassword     string `json:"newPassword" validate:"required,min=6"`
}
