package types

// ConfigModel model
type ConfigModel struct {
	Secrets          []string         `json:"secrets,omitempty"`
	BasePath         string           `json:"base_path,omitempty"`
	CatalogueHostAPI string           `json:"catalogueHostApi"`
	TMSEndpoint      string           `json:"tmsEndpoint"`
	Database         EngineDB         `json:"db,omitempty"`
	Server           HTTPServerConfig `json:"http"`
}

type DBConfig struct {
	Dialect        string `json:"dialect"`
	Database       string `json:"database"`
	Host           string `json:"host"`
	Port           int    `json:"port"`
	User           string `json:"user"`
	Password       string `json:"password"`
	MaxConnections int    `json:"max_conns"`
}

// EngineDB model
type EngineDB struct {
	Writer *DBConfig `json:"writer,omitempty"`
	Reader *DBConfig `json:"reader,omitempty"`
}

type HTTPServerConfig struct {
	BasePath   *string `json:"basePath"`
	ServerHost string  `json:"serverHost"`
	ServerPort int     `json:"serverPort"`
}
