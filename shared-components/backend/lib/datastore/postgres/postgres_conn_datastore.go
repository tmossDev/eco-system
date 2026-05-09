package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"tmossDev.github.com/eco-system/shared-components/backend/lib/config/types"
	"tmossDev.github.com/eco-system/shared-components/backend/lib/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/lib/logger"
)

type PostgresDataStore struct {
	WriterDB *sql.DB
	ReaderDB *sql.DB
	config   types.EngineDB
}

func NewPostgresDataStore(config types.EngineDB) *PostgresDataStore {
	return &PostgresDataStore{WriterDB: nil, ReaderDB: nil, config: config}
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

	r.WriterDB = writerDB
	r.ReaderDB = readerDB

	return nil
}

func createConnection(dbConfig *types.DBConfig) (*sql.DB, error) {
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

	// ✅ Ping with timeout (better practice)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func (r *PostgresDataStore) Close() (error, error) {
	wErr := r.WriterDB.Close()
	rErr := r.ReaderDB.Close()
	return wErr, rErr
}

func (r *PostgresDataStore) Ping() error {
	if err := r.WriterDB.Ping(); err != nil {
		return err
	}
	if err := r.ReaderDB.Ping(); err != nil {
		return err
	}
	return nil
}

func (r *PostgresDataStore) CreateWriterTransaction() (*sql.Tx, error) {
	return r.WriterDB.Begin()
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
