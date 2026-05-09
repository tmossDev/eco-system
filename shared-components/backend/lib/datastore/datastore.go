package datastore

type DataStore interface {
	Connect() error
	Close() error
	Ping() error
}
