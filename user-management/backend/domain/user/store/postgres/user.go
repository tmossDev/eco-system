package postgres

import (
	"database/sql"
	"errors"
	"time"

	"tmossDev.github.com/eco-system/shared-components/backend/package/datastore"
	"tmossDev.github.com/eco-system/shared-components/backend/package/datastore/flows"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
	"tmossDev.github.com/eco-system/shared-components/backend/package/user/model"
	"tmossDev.github.com/eco-system/shared-components/backend/package/utils"
	"tmossDev.github.com/eco-system/user-management/backend/domain/user/repository"
)

type PostgresUserRepository struct {
	store datastore.DataStore
}

const MySystemAutoID = 1

func NewPostgresUserRepository(store datastore.DataStore) repository.UserRepository {
	return &PostgresUserRepository{
		store: store,
	}
}

func (repo *PostgresUserRepository) Shutdown() {
	err := repo.store.Close()
	if err != nil {
		logger.Errorf("Unabled to close user repo: %s", err.Error())
	}
}

func (repo *PostgresUserRepository) mapStatementToUser(row *sql.Row) (*model.UserResponse, error) {
	var user model.UserResponse

	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.HashedPassword,
		&user.RoleID,
		&user.CreatedUser,
		&user.CreatedAt,
		&user.UpdatedUser,
		&user.UpdatedAt,
		&user.DeletedUser,
		&user.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debugf("No result back for user: %s", err.Error())
			return nil, types.NewNoTFoundOrNoRecordError()
		}

		logger.Errorf("Unabled to marshal user response: %s", err.Error())
		return nil, err
	}

	return &user, nil
}

func (repo *PostgresUserRepository) GetByEmail(requestId string, email string) (*model.UserResponse, error) {
	stmt, err := flows.GetReaderStatement("GetByEmail", GetByEmail, repo.store)
	if err != nil {
		return nil, err
	}

	logger.Debugf("Running query '%s' with parameter '%s'", GetByEmail, email)

	user, err := repo.mapStatementToUser(stmt.QueryRow(email))
	if err != nil {
		utils.LogExecutingError("GetByEmail", err, requestId)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.NewNoTFoundOrNoRecordError()
		}

		return nil, types.NewInternalServerError()
	}

	return user, nil
}

func (repo *PostgresUserRepository) GetByID(userId uint64) (*model.UserResponse, error) {
	stmt, err := flows.GetReaderStatement("GetByID", GetByID, repo.store)
	if err != nil {
		return nil, err
	}

	logger.Debugf("Running query '%s' with parameter '%d'", GetByID, userId)

	user, err := repo.mapStatementToUser(stmt.QueryRow(userId))
	if err != nil {
		utils.LogExecutingError("GetByID", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.NewNoTFoundOrNoRecordError()
		}

		return nil, types.NewInternalServerError()
	}

	return user, nil
}

func (repo *PostgresUserRepository) RegisterUser(firstName string, lastName string, email string, password string, roleId uint64) (*model.UserResponse, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		logger.Errorf("Error creating hashPassword: %s", err.Error())
		return nil, types.NewInternalServerError()
	}

	insertedAt := utils.GetCurrentDateFormatedForInsertingIntoDB(time.Now())

	logger.Debugf(
		"Running query '%s' with parameter '%s', '%s', '%s', '%v', '%d', '%d' and '%s'",
		RegisterUser,
		firstName,
		lastName,
		email,
		hashedPassword,
		roleId,
		MySystemAutoID,
		insertedAt,
	)

	insertedID, err := flows.PerformEdit(
		"RegisterUser",
		RegisterUser,
		repo.store,
		firstName,
		lastName,
		email,
		hashedPassword,
		roleId,
		MySystemAutoID,
		insertedAt,
	)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:          uint64(insertedID),
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		RoleID:      roleId,
		CreatedUser: MySystemAutoID,
		CreatedAt:   insertedAt,
	}, nil
}

func (repo *PostgresUserRepository) Update(userId uint64, firstName string, lastName string, updatingUserId uint64) (*model.UserResponse, error) {

	logger.Debugf(
		"Running query '%s' with parameter '%s', '%s', '%d' and '%d'",
		UpdateUser,
		firstName,
		lastName,
		updatingUserId,
		userId,
	)

	_, err := flows.PerformEdit(
		"UpdateUser",
		UpdateUser,
		repo.store,
		firstName,
		lastName,
		updatingUserId,
		userId,
	)
	if err != nil {
		return nil, err
	}

	return repo.GetByID(userId)
}

func (repo *PostgresUserRepository) ResetPassword(userId uint64, newPassword string, updatingUserId uint64) (*model.UserResponse, error) {
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		logger.Errorf("Error hashing password: %s", err.Error())
		return nil, types.NewInternalServerError()
	}

	logger.Debugf(
		"Running query '%s' with parameter '%v', '%d' and '%d'",
		ResetPassword,
		hashedPassword,
		updatingUserId,
		userId,
	)

	_, err = flows.PerformEdit(
		"ResetPassword",
		ResetPassword,
		repo.store,
		hashedPassword,
		updatingUserId,
		userId,
	)
	if err != nil {
		return nil, err
	}

	return repo.GetByID(userId)
}

func (repo *PostgresUserRepository) ResetEmail(userId uint64, newEmail string, updatingUserId uint64) (*model.UserResponse, error) {
	logger.Debugf(
		"Running query '%s' with parameter '%s', '%d' and '%d'",
		ResetEmail,
		newEmail,
		updatingUserId,
		userId,
	)

	_, err := flows.PerformEdit(
		"ResetEmail",
		ResetEmail,
		repo.store,
		newEmail,
		updatingUserId,
		userId,
	)
	if err != nil {
		return nil, err
	}

	return repo.GetByID(userId)
}
