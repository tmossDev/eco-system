package types

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	configTypes "tmossDev.github.com/eco-system/shared-components/backend/package/config/types"
)

type MySql interface {
	Ping() error
	Close() error
	GetReader() error
	GetWriter() error
	createConnection(dbConfig configTypes.DBConfig) (*sql.DB, error)
}

type DbConnection struct {
	readerDB *sql.DB
	writerDB *sql.DB
}

func NewDbConnection(writerDB, readerDB *sql.DB) (*DbConnection, error) {

	return &DbConnection{readerDB: readerDB, writerDB: writerDB}, nil
}

func createConnection(dbConfig configTypes.DBConfig) (*sql.DB, error) {

	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
	)

	db, err := sql.Open(dbConfig.Dialect, connString)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(3)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DbConnection) GetReader() *sql.DB {
	return db.readerDB
}

func (db *DbConnection) GetWriter() *sql.DB {
	return db.writerDB
}

func (db *DbConnection) Close() error {
	var err = db.readerDB.Close()
	if err != nil {
		return err
	}

	return db.writerDB.Close()
}

func (db *DbConnection) Ping() error {
	err := db.readerDB.Ping()
	if err != nil {
		return err
	}
	return db.writerDB.Ping()
}
