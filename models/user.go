package models

type UserRole string
type UserGender string

const (
	RoleUser     UserRole   = "user"
	RoleAdmin    UserRole   = "admin"
	GenderFemale UserGender = "female"
	GenderMale   UserGender = "male"
	GenderOther  UserGender = "other"
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
	Username    string     `json:"username" validate:"required,min=3,max=50"`
	FullName    string     `json:"fullName" validate:"required"`
	Password    string     `json:"password" validate:"required,min=6"`
	Gender      UserGender `json:"gender" validate:"required,oneof=female male other prefer_not_to_say"`
	Avatar      string     `json:"avatar"`
	DateOfBirth string     `json:"dateOfBirth" validate:"required,datetime=2006-01-02"`
	PhoneNumber string     `json:"phoneNumber" validate:"omitempty,e164"`
	Role        UserRole   `json:"role"`
}

type UserDetail struct {
	ID          int        `json:"id" validate:"required"`
	Email       string     `json:"email" validate:"required,email"`
	Username    string     `json:"username" validate:"required"`
	FullName    string     `json:"fullName" validate:"required"`
	Gender      UserGender `json:"gender" validate:"required"`
	Avatar      string     `json:"avatar,omitempty"`
	DateOfBirth string     `json:"dateOfBirth" validate:"required,datetime=2006-01-02"`
	PhoneNumber string     `json:"phoneNumber,omitempty" validate:"omitempty,e164"`
	Role        UserRole   `json:"role" validate:"required"`
}

type UserToken struct {
	ID   int      `json:"id" validate:"required"`
	Role UserRole `json:"role" validate:"required"`
}

type UpdateUserResponse struct {
	Message string     `json:"message" validate:"required"`
	User    UserDetail `json:"user" validate:"required"`
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

type CreateUserResponse struct {
	Message string     `json:"message"`
	User    UserDetail `json:"user"`
}
