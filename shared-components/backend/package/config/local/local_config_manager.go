package local

import (
	"os"
	"strconv"

	"tmossDev.github.com/eco-system/shared-components/backend/package/config"
	"tmossDev.github.com/eco-system/shared-components/backend/package/config/types"
)

type Config struct {
}

func NewLocalConfigManager() config.Config {
	return &Config{}
}

func (sh *Config) GetConfig(_ string) (*types.ConfigModel, error) {
	dbPort, err := strconv.Atoi(getenv("DB_PORT", "5432"))
	if err != nil {
		dbPort = 5432
	}

	maxConnections, err := strconv.Atoi(getenv("DB_MAX_CONNECTIONS", "15"))
	if err != nil {
		maxConnections = 15
	}

	db := &types.DBConfig{
		Dialect:        getenv("DB_DIALECT", "postgresql"),
		Database:       getenv("DB_NAME", "ecoDB"),
		Host:           getenv("DB_HOST", "db"),
		Port:           dbPort,
		User:           getenv("DB_USER", "postgres"),
		Password:       getenv("DB_PASSWORD", "very_secure_password"),
		MaxConnections: maxConnections,
	}
	basePath := "/"
	return &types.ConfigModel{
		Server: types.HTTPServerConfig{
			ServerPort: 3001,
			BasePath:   &basePath,
		},

		Database: types.EngineDB{
			Writer: db,
			Reader: db,
		},
		BasePath: basePath,
	}, nil
}

func getenv(name string, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}

	return value
}
