package local

import (
	"tmossDev.github.com/eco-system/shared-components/backend/package/config"
	"tmossDev.github.com/eco-system/shared-components/backend/package/config/types"
)

type Config struct {
}

func NewLocalConfigManager() config.Config {
	return &Config{}
}

func (sh *Config) GetConfig(_ string) (*types.ConfigModel, error) {
	db := &types.DBConfig{
		Dialect:        "postgresql",
		Database:       "ecoDB",
		Host:           "db",
		Port:           5432,
		User:           "postgres",
		Password:       "very_secure_password",
		MaxConnections: 15,
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
