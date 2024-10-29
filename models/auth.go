package models

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

type AccessTokenResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int64  `json:"expiresIn"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Password string `json:"password" validate:"required,min=6"`
}
