package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"tmossDev.github.com/eco-system/shared-components/backend/package/datastore/types"

	configTypes "tmossDev.github.com/eco-system/shared-components/backend/package/config/types"
	"tmossDev.github.com/eco-system/shared-components/backend/package/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
)

type PostgresDataStore struct {
	conn   *types.DbConnection
	config configTypes.EngineDB
}

func (r *PostgresDataStore) GetConnection() *types.DbConnection {
	return r.conn
}

func NewPostgresDataStore(config configTypes.EngineDB) *PostgresDataStore {
	return &PostgresDataStore{conn: nil, config: config}
}

func (r *PostgresDataStore) Connect() error {
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
		dbConfig.Port = 5432
	}

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
	)

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(1)

	if dbConfig.MaxConnections > 0 {
		db.SetMaxOpenConns(dbConfig.MaxConnections)
	}

	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func (r *PostgresDataStore) Close() error {
	err := r.conn.GetWriter().Close()
	if err != nil {
		return err
	}
	err = r.conn.GetReader().Close()
	return err
}

func (r *PostgresDataStore) Ping() error {
	if err := r.conn.GetWriter().Ping(); err != nil {
		return err
	}
	if err := r.conn.GetReader().Ping(); err != nil {
		return err
	}
	return nil
}

func (r *PostgresDataStore) CreateWriterTransaction() (*sql.Tx, error) {
	return r.conn.GetWriter().Begin()
}

func (r *PostgresDataStore) NewSqlContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func (r *PostgresDataStore) RollbackAndJoinErrorIfAny(tx *sql.Tx) {
	if tx == nil {
		return
	}

	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		logger.Errorf(constants.DefaultRequestId, "unable to rollback sql transaction: %s", err.Error())
	}
}

func (r *PostgresDataStore) CloseStatement(stmt *sql.Stmt) {
	if stmt == nil {
		return
	}

	if err := stmt.Close(); err != nil {
		logger.Errorf(constants.DefaultRequestId, "unable to close sql statement: %s", err.Error())
	}
}

func (r *PostgresDataStore) CloseRows(rows *sql.Rows) {
	if rows == nil {
		return
	}

	if err := rows.Close(); err != nil {
		logger.Warnf(constants.DefaultRequestId, "unable to close rows: %s", err.Error())
	}
}
