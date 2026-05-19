package repository

import (
	"tmossDev.github.com/eco-system/user-management/backend/package/user/model"
)

type UserRepository interface {
	GetByID(userId uint64) (*model.UserResponse, error)
	GetByEmail(requestId string, email string) (*model.UserResponse, error)
	RegisterUser(firstName string, lastName string, email string, password string, roleId uint64) (*model.UserResponse, error)
	Update(userId uint64, firstName string, lastName string, updatingUserId uint64) (*model.UserResponse, error)
	ResetPassword(userId uint64, newPassword string, updatingUserId uint64) (*model.UserResponse, error)
	ResetEmail(userId uint64, newEmail string, updatingUserId uint64) (*model.UserResponse, error)
	Shutdown()
}
