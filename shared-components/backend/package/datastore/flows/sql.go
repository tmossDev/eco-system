package flows

import (
	"database/sql"

	"tmossDev.github.com/eco-system/shared-components/backend/package/datastore"
	libTypes "tmossDev.github.com/eco-system/shared-components/backend/package/types"
	"tmossDev.github.com/eco-system/shared-components/backend/package/utils"
)

func GetReaderStatement(queryName string, query string, conn datastore.DataStore) (*sql.Stmt, error) {
	stmt, err := conn.GetConnection().GetReader().Prepare(query)
	if err != nil {
		utils.LogPreparingError(queryName, err)
		return nil, libTypes.NewInternalServerError()
	}

	return stmt, nil
}

func PerformEdit(queryName string, query string, conn datastore.DataStore, args ...any) (int64, error) {
	tx, err := conn.GetConnection().GetWriter().Begin()
	if err != nil {
		utils.LogPreparingError(queryName, err)
		return -1, libTypes.NewInternalServerError()
	}

	preparedStmt, err := tx.Prepare(query)
	if err != nil {
		utils.LogPreparingError(queryName, err)
		tx.Rollback()
		return -1, libTypes.NewInternalServerError()
	}
	defer preparedStmt.Close()

	result, err := preparedStmt.Exec(args...)
	if err != nil {
		utils.LogExecutingError(queryName, err)
		tx.Rollback()
		return -1, libTypes.NewInternalServerError()
	}

	err = tx.Commit()
	if err != nil {
		utils.LogCommitError(queryName, err)
		return -1, libTypes.NewInternalServerError()
	}

	id, _ := result.LastInsertId()

	return id, nil
}
