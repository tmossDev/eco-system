package model

type ChangePasswordRequest struct {
	Password        string `json:"password" validate:"required,gt=0"`
	ConfirmPassowrd string `json:"confirm_password" validate:"required,gt=0"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,gt=0,lte=225"`
	Password string `json:"password" validate:"required,gt=0"`
}

type ChangeEmailRequest struct {
	Email        string `json:"email" validate:"required,gt=0,lte=225"`
	ConfirmEmail string `json:"confirm_email" validate:"required,gt=0,lte=225"`
}

type UserUpdateRequest struct {
	FirstName string `json:"first_name" validate:"required,gt=0,lte=50"`
	LastName  string `json:"last_name" validate:"required,gt=0,lte=50"`
}

type UserRequest struct {
	UserUpdateRequest
	Email           string `json:"email" validate:"required,gt=0,lte=225"`
	Password        string `json:"password" validate:"required,gt=0"`
	ConfirmPassword string `json:"confirm_password" validate:"required,gt=0"`
}
