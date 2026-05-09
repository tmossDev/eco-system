package config

import (
	"tmossDev.github.com/eco-system/shared-components/backend/lib/config/types"
)

type Config interface {
	GetConfig(configName string) (*types.ConfigModel, error)
}
