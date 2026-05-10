package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	configTypes "tmossDev.github.com/eco-system/shared-components/backend/package/config/types"
	"tmossDev.github.com/eco-system/shared-components/backend/package/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/datastore/types"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
)

type MySqlDataStore struct {
	conn   *types.DbConnection
	config configTypes.EngineDB
}

func NewMysqlDataStore(config configTypes.EngineDB) *MySqlDataStore {
	return &MySqlDataStore{conn: nil, config: config}
}

func (r *MySqlDataStore) Connect() error {
	writerDB, err := createConnection(r.config.Writer)
	if err != nil {
		return err
	}
	readerDB, err := createConnection(r.config.Reader)
	if err != nil {
		writerDB.Close()
		return err
	}

	r.conn, err = types.NewDbConnection(writerDB, readerDB)
	if err != nil {
		readerDB.Close()
		writerDB.Close()
		return err
	}

	return nil
}

func createConnection(dbConfig *configTypes.DBConfig) (*sql.DB, error) {
	if dbConfig.Port == 0 {
		dbConfig.Port = 3306
	}
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	db, err := sql.Open(dbConfig.Dialect, connString)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(1)
	if dbConfig.MaxConnections > 0 {
		db.SetMaxOpenConns(dbConfig.MaxConnections)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (r *MySqlDataStore) Close() (error, error) {
	wErr := r.conn.GetWriter().Close()
	rErr := r.conn.GetReader().Close()
	return wErr, rErr
}

func (r *MySqlDataStore) Ping() error {
	err := r.conn.GetWriter().Ping()
	if err != nil {
		return err
	}
	err = r.conn.GetReader().Ping()
	if err != nil {
		return err
	}
	return nil
}

func (r *MySqlDataStore) CreateWriterTransaction() (*sql.Tx, error) {
	return r.conn.GetWriter().Begin()
}

func (r *MySqlDataStore) NewSqlContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx, cancel
}

func (r *MySqlDataStore) RollbackAndJoinErrorIfAny(tx *sql.Tx) {
	rollBackErr := tx.Rollback()
	if rollBackErr != nil {
		logger.Errorf(constants.DefaultRequestId, "unabled to rolw back sql trasnaction reason: %s", rollBackErr.Error())
	}
}

func (r *MySqlDataStore) CloseStatement(stmt *sql.Stmt) {
	err := stmt.Close()
	if err != nil {
		logger.Errorf(constants.DefaultRequestId, "unabled to close sql statement reason: %s", err.Error())
	}
}

func (r *MySqlDataStore) CloseRows(row *sql.Rows) {
	err := row.Close()
	if err != nil {
		logger.Warnf(constants.DefaultRequestId, "unable to close row(s) reason: %s", err.Error())
	}
}
