package repository

import (
	"tmossDev.github.com/eco-system/shared-components/backend/package/user/model"
)

type UserRepository interface {
	GetByID(userId uint64) (*model.UserResponse, error)
	GetByEmail(requestId string, email string) (*model.UserResponse, error)
	RegisterUser(user *model.UserResponse) error
	Update(user model.UserResponse) error
	ResetPassword(user model.UserResponse) error
	ResetEmail(user model.UserResponse) error
	Shutdown()
}
