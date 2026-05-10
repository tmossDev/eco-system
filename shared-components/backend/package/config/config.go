package config

import (
	"tmossDev.github.com/eco-system/shared-components/backend/package/config/types"
)

type Config interface {
	GetConfig(configName string) (*types.ConfigModel, error)
}
