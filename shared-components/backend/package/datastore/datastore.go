package datastore

import "tmossDev.github.com/eco-system/shared-components/backend/package/datastore/types"

type DataStore interface {
	Connect() error
	Close() error
	Ping() error
	GetConnection() *types.DbConnection
}
