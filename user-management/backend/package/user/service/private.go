package service

import (
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
	"tmossDev.github.com/eco-system/shared-components/backend/package/validator"
	"tmossDev.github.com/eco-system/user-management/backend/package/user/model"
	"tmossDev.github.com/eco-system/user-management/backend/package/user/repository"
)

type PrivateUserService interface {
	UpdateUserInfo(userId uint64, body string, updatingUserId uint64) (*model.UserResponse, error)
	UpdateUserPassword(userId uint64, body string, updatingUserId uint64) (*model.UserResponse, error)
	Shutdown()
}

type PrivateUserServiceImpl struct {
	validator validator.Validator
	userRepo  repository.UserRepository
}

func (p *PrivateUserServiceImpl) UpdateUserInfo(userId uint64, body string, updatingUserId uint64) (*model.UserResponse, error) {
	var updateUserRequest model.UserUpdateRequest
	err := p.validator.MarshalAndValidateREQ(body, &updateUserRequest)
	if err != nil {
		return nil, err
	}

	user, err := p.userRepo.Update(userId, updateUserRequest.FirstName, updateUserRequest.LastName, updatingUserId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *PrivateUserServiceImpl) UpdateUserPassword(userId uint64, body string, updatingUserId uint64) (*model.UserResponse, error) {
	var updateUserRequest model.ChangePasswordRequest
	err := p.validator.MarshalAndValidateREQ(body, &updateUserRequest)
	if err != nil {
		return nil, err
	}

	if updateUserRequest.ConfirmPassowrd != updateUserRequest.Password {
		return nil, types.NewInvalidInputError()
	}

	user, err := p.userRepo.ResetPassword(userId, updateUserRequest.Password, updatingUserId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *PrivateUserServiceImpl) Shutdown() {
	p.userRepo.Shutdown()
}

func NewPrivateService(validator validator.Validator, userRepo repository.UserRepository) PrivateUserService {
	return &PrivateUserServiceImpl{
		validator: validator,
		userRepo:  userRepo,
	}
}
