package service

import (
	"strconv"

	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
	"tmossDev.github.com/eco-system/shared-components/backend/package/utils"
	"tmossDev.github.com/eco-system/shared-components/backend/package/validator"
	"tmossDev.github.com/eco-system/user-management/backend/package/user/constants"
	"tmossDev.github.com/eco-system/user-management/backend/package/user/model"
	"tmossDev.github.com/eco-system/user-management/backend/package/user/repository"
)

type PublicUserService interface {
	IsAuthenticated(jwt string) error
	IsAuthorized(jwt string, page string) error
	User(jwt string) (*model.UserResponse, error)
	Logout(jwt string) error
	Register(body string) (*model.LoginResponse, error)
	Login(body string) (*model.LoginResponse, error)
	Shutdown()
}

type PublicUserServiceImpl struct {
	validator validator.Validator
	userRepo  repository.UserRepository
}

func (auth *PublicUserServiceImpl) Shutdown() {
	auth.userRepo.Shutdown()
	auth = nil
}

func NewPublicService(validator validator.Validator, userReo repository.UserRepository) PublicUserService {
	return &PublicUserServiceImpl{
		validator: validator,
		userRepo:  userReo,
	}
}

func (auth *PublicUserServiceImpl) generateLoginResponseFromUser(user model.UserResponse) (*model.LoginResponse, error) {

	token, expireAt, err := utils.GenerateJwt(utils.UintToString(user.ID), constants.PASSWORD_SECRET_HASHING_KEY)
	if err != nil {
		logger.Infof("Cant generate Jwt for userId = '%d' due to '%s'", strconv.FormatUint(user.ID, 10), err.Error())
		return nil, types.NewInternalServerError()
	}
	logger.Debugf("User logged in '%d' with token '%s'", strconv.FormatUint(user.ID, 10), token)

	return &model.LoginResponse{
		Jwt:      token,
		ExpireAt: expireAt,
	}, nil
}

func (auth *PublicUserServiceImpl) checkJwt(jwt string) (string, error) {
	Id, err := utils.ParseJwt(jwt, constants.PASSWORD_SECRET_HASHING_KEY)

	if err != nil {
		logger.Infof("Token Unabled to Parse: '%s'", err.Error())
		return "", err
	}

	logger.Debugf("Token Parsed with ID: '%s'", Id)

	return Id, nil
}

func (auth *PublicUserServiceImpl) IsAuthenticated(jwt string) error {
	if _, err := auth.checkJwt(jwt); err != nil {
		return err
	}

	return nil
}

func (auth *PublicUserServiceImpl) IsAuthorized(jwt string, page string) error {
	_, err := auth.checkJwt(jwt)
	if err != nil {
		return err
	}

	// check page permission

	return nil
}

func (auth PublicUserServiceImpl) Login(body string) (*model.LoginResponse, error) {
	var loginRequest model.LoginRequest
	err := auth.validator.MarshalAndValidateREQ(body, &loginRequest)
	if err != nil {
		return nil, err
	}

	user, err := auth.userRepo.GetByEmail(loginRequest.Username)
	if err != nil {
		return nil, err
	}

	if !utils.ComparePassword(user.HashedPassword, loginRequest.Password) {
		logger.Infof("Failed login attempt for '%s' with '%s'", loginRequest.Username, loginRequest.Password)
		return nil, types.NewUnauthorizedError()
	}

	return auth.generateLoginResponseFromUser(*user)

}

func (auth *PublicUserServiceImpl) Logout(jwt string) error {
	userId, err := utils.GetIssuerFromJwt(jwt, constants.PASSWORD_SECRET_HASHING_KEY)
	if err != nil {
		logger.Errorf("Unabled to logout token: '%s'", jwt)
	}
	logger.Debugf("Logged out user with Id = %d and token: '%s'", userId, jwt)
	return nil
}

func (auth *PublicUserServiceImpl) Register(body string) (*model.LoginResponse, error) {
	var registerRequest model.UserRequest
	err := auth.validator.MarshalAndValidateREQ(body, &registerRequest)
	if err != nil {
		return nil, err
	}

	if registerRequest.ConfirmPassword != registerRequest.Password {
		return nil, types.NewInvalidInputError()
	}

	user, err := auth.userRepo.RegisterUser(registerRequest.FirstName, registerRequest.LastName, registerRequest.Email, registerRequest.Password, constants.CUSTOMER_ROLE_ID)
	if err != nil {
		return nil, err
	}

	return auth.generateLoginResponseFromUser(*user)
}

func (auth *PublicUserServiceImpl) User(jwt string) (*model.UserResponse, error) {
	issuer, err := utils.GetIssuerFromJwt(jwt, constants.PASSWORD_SECRET_HASHING_KEY)
	if err != nil {
		logger.Errorf("Cant get issuer from jwt '%s'", issuer)
		return nil, types.NewInternalServerError()
	}
	userId, err := strconv.ParseUint(issuer, 10, 64)
	if err != nil {
		logger.Errorf("Cant parse as Uint: '%s'", issuer)
		return nil, types.NewInternalServerError()
	}
	user, err := auth.userRepo.GetByID(userId)
	if err != nil {
		logger.Errorf("No user from token: '%s'", jwt)
		return nil, types.NewInternalServerError()
	}
	return user, nil
}
