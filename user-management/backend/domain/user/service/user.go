package service

import (
	"errors"
	"strconv"
	"time"

	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
	"tmossDev.github.com/eco-system/shared-components/backend/package/user/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/user/model"
	"tmossDev.github.com/eco-system/shared-components/backend/package/utils"
	"tmossDev.github.com/eco-system/shared-components/backend/package/validator"
	"tmossDev.github.com/eco-system/user-management/backend/domain/user/repository"
)

type UserService interface {
	IsAuthenticated(jwt string) error
	IsAuthorized(jwt string, page string) error
	User(jwt string) (*model.UserResponse, error)
	Logout(jwt string) error
	Register(body string) (*model.LoginResponse, error)
	Login(requestId string, body string) (*model.LoginResponse, error)
	UpdateUserInfo(userId uint64, body string, updatingUserId uint64) (*model.UserResponse, error)
	UpdateUserPassword(userId uint64, body string, updatingUserId uint64) (*model.UserResponse, error)
	Shutdown()
}

type UserServiceImpl struct {
	validator validator.Validator
	userRepo  repository.UserRepository
}

func (service *UserServiceImpl) Shutdown() {
	service.userRepo.Shutdown()
}

func NewUserService(validator validator.Validator, userRepo repository.UserRepository) UserService {
	return &UserServiceImpl{
		validator: validator,
		userRepo:  userRepo,
	}
}

func (service *UserServiceImpl) generateLoginResponseFromUser(user model.UserResponse) (*model.LoginResponse, error) {

	token, expireAt, err := utils.GenerateJwt(utils.UintToString(user.ID), constants.PASSWORD_SECRET_HASHING_KEY)
	if err != nil {
		logger.Infof("Cant generate Jwt for userId = '%d' due to '%s'", strconv.FormatUint(user.ID, 10), err.Error())
		return nil, types.NewInternalServerError()
	}
	logger.Debugf("User logged in '%d' with token '%s'", strconv.FormatUint(user.ID, 10), token)

	return &model.LoginResponse{
		Jwt:         token,
		AccessToken: token,
		ExpireAt:    expireAt,
		User: model.AuthUserResponse{
			ID:    utils.UintToString(user.ID),
			Name:  user.FirstName + " " + user.LastName,
			Email: user.Email,
			Role:  strconv.FormatUint(user.RoleID, 10),
		},
	}, nil
}

func (service *UserServiceImpl) checkJwt(jwt string) (string, error) {
	Id, err := utils.ParseJwt(jwt, constants.PASSWORD_SECRET_HASHING_KEY)

	if err != nil {
		logger.Infof("Token Unabled to Parse: '%s'", err.Error())
		return "", err
	}

	logger.Debugf("Token Parsed with ID: '%s'", Id)

	return Id, nil
}

func (service *UserServiceImpl) IsAuthenticated(jwt string) error {
	if _, err := service.checkJwt(jwt); err != nil {
		return err
	}

	return nil
}

func (service *UserServiceImpl) IsAuthorized(jwt string, page string) error {
	_, err := service.checkJwt(jwt)
	if err != nil {
		return err
	}

	// check page permission

	return nil
}

func (service *UserServiceImpl) Login(requestId string, body string) (*model.LoginResponse, error) {
	var loginRequest model.LoginRequest
	err := service.validator.MarshalAndValidateREQ(body, &loginRequest)
	if err != nil {
		return nil, err
	}

	username := loginRequest.Username
	if username == "" {
		username = loginRequest.Email
	}
	if username == "" {
		logger.Infof(requestId, "Login failed: no username or email supplied")
		return nil, types.NewUnauthorizedError()
	}

	user, err := service.userRepo.GetByEmail(requestId, username)
	if err != nil {
		var socketErr *types.SocketError
		if errors.As(err, &socketErr) {
			logger.Infof(requestId, "Login failed for '%s': %s", username, socketErr.Error())
			return nil, types.NewUnauthorizedError()
		}

		logger.Errorf(requestId, "Login lookup failed for '%s': %s", username, err.Error())
		return nil, err
	}

	if !utils.ComparePassword(user.HashedPassword, loginRequest.Password) {
		logger.Infof(requestId, "Failed login attempt for '%s': invalid password", username)
		return nil, types.NewUnauthorizedError()
	}

	return service.generateLoginResponseFromUser(*user)

}

func (service *UserServiceImpl) Logout(jwt string) error {
	userId, err := utils.GetIssuerFromJwt(jwt, constants.PASSWORD_SECRET_HASHING_KEY)
	if err != nil {
		logger.Errorf("Unabled to logout token: '%s'", jwt)
	}
	logger.Debugf("Logged out user with Id = %d and token: '%s'", userId, jwt)
	return nil
}

func (service *UserServiceImpl) Register(body string) (*model.LoginResponse, error) {
	var registerRequest model.UserRequest
	err := service.validator.MarshalAndValidateREQ(body, &registerRequest)
	if err != nil {
		return nil, err
	}

	if registerRequest.ConfirmPassword != registerRequest.Password {
		return nil, types.NewInvalidInputError()
	}

	hashedPassword, err := hashPassword(registerRequest.Password)
	if err != nil {
		return nil, err
	}

	user := model.UserResponse{
		FirstName:      registerRequest.FirstName,
		LastName:       registerRequest.LastName,
		Email:          registerRequest.Email,
		HashedPassword: hashedPassword,
		RoleID:         constants.CUSTOMER_ROLE_ID,
		CreatedUser:    MySystemAutoID,
		CreatedAt:      utils.GetCurrentDateFormatedForInsertingIntoDB(time.Now()),
	}
	if err := service.userRepo.RegisterUser(&user); err != nil {
		return nil, err
	}

	return service.generateLoginResponseFromUser(user)
}

func (service *UserServiceImpl) User(jwt string) (*model.UserResponse, error) {
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
	user, err := service.userRepo.GetByID(userId)
	if err != nil {
		logger.Errorf("No user from token: '%s'", jwt)
		return nil, types.NewInternalServerError()
	}
	return user, nil
}

func (service *UserServiceImpl) UpdateUserInfo(userId uint64, body string, updatingUserId uint64) (*model.UserResponse, error) {
	var updateUserRequest model.UserUpdateRequest
	err := service.validator.MarshalAndValidateREQ(body, &updateUserRequest)
	if err != nil {
		return nil, err
	}

	user, err := service.userRepo.GetByID(userId)
	if err != nil {
		return nil, err
	}

	user.FirstName = updateUserRequest.FirstName
	user.LastName = updateUserRequest.LastName
	user.UpdatedUser = &updatingUserId
	if err := service.userRepo.Update(*user); err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserServiceImpl) UpdateUserPassword(userId uint64, body string, updatingUserId uint64) (*model.UserResponse, error) {
	var updateUserRequest model.ChangePasswordRequest
	err := service.validator.MarshalAndValidateREQ(body, &updateUserRequest)
	if err != nil {
		return nil, err
	}

	if updateUserRequest.ConfirmPassowrd != updateUserRequest.Password {
		return nil, types.NewInvalidInputError()
	}

	hashedPassword, err := hashPassword(updateUserRequest.Password)
	if err != nil {
		return nil, err
	}

	user, err := service.userRepo.GetByID(userId)
	if err != nil {
		return nil, err
	}

	user.HashedPassword = hashedPassword
	user.UpdatedUser = &updatingUserId
	if err := service.userRepo.ResetPassword(*user); err != nil {
		return nil, err
	}

	return user, nil
}

const MySystemAutoID = 1

func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		logger.Errorf("Error hashing password: %s", err.Error())
		return nil, types.NewInternalServerError()
	}

	return hashedPassword, nil
}
