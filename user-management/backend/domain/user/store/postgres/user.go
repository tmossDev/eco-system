package postgres

import (
	"database/sql"
	"errors"

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

func (repo *PostgresUserRepository) RegisterUser(user *model.UserResponse) error {
	logger.Debugf(
		"Running query '%s' with parameter '%s', '%s', '%s', '%v', '%d', '%d' and '%s'",
		RegisterUser,
		user.FirstName,
		user.LastName,
		user.Email,
		user.HashedPassword,
		user.RoleID,
		user.CreatedUser,
		user.CreatedAt,
	)

	insertedID, err := flows.PerformEdit(
		"RegisterUser",
		RegisterUser,
		repo.store,
		user.FirstName,
		user.LastName,
		user.Email,
		user.HashedPassword,
		user.RoleID,
		user.CreatedUser,
		user.CreatedAt,
	)
	if err != nil {
		return err
	}

	user.ID = uint64(insertedID)
	return nil
}

func (repo *PostgresUserRepository) Update(user model.UserResponse) error {
	logger.Debugf(
		"Running query '%s' with parameter '%s', '%s', '%d' and '%d'",
		UpdateUser,
		user.FirstName,
		user.LastName,
		user.UpdatedUser,
		user.ID,
	)

	_, err := flows.PerformEdit(
		"UpdateUser",
		UpdateUser,
		repo.store,
		user.FirstName,
		user.LastName,
		user.UpdatedUser,
		user.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PostgresUserRepository) ResetPassword(user model.UserResponse) error {
	logger.Debugf(
		"Running query '%s' with parameter '%v', '%d' and '%d'",
		ResetPassword,
		user.HashedPassword,
		user.UpdatedUser,
		user.ID,
	)

	_, err := flows.PerformEdit(
		"ResetPassword",
		ResetPassword,
		repo.store,
		user.HashedPassword,
		user.UpdatedUser,
		user.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PostgresUserRepository) ResetEmail(user model.UserResponse) error {
	logger.Debugf(
		"Running query '%s' with parameter '%s', '%d' and '%d'",
		ResetEmail,
		user.Email,
		user.UpdatedUser,
		user.ID,
	)

	_, err := flows.PerformEdit(
		"ResetEmail",
		ResetEmail,
		repo.store,
		user.Email,
		user.UpdatedUser,
		user.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
