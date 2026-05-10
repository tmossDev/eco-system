package utils

import (
	"fmt"

	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
)

func LogPreparingError(queryName string, err error) {
	logger.Error(fmt.Sprintf("Error preparing '%s' query: %s", queryName, err.Error()))
}

func LogExecutingError(queryName string, err error) {
	logger.Error(fmt.Sprintf("Error executing '%s' query: %s", queryName, err.Error()))
}

func LogBeginingTnxError(queryName string, err error) {
	logger.Error(fmt.Sprintf("Error begining transaction '%s' query: %s", queryName, err.Error()))
}

func LogCommitError(queryName string, err error) {
	logger.Error(fmt.Sprintf("Error commiting transaction '%s' query: %s", queryName, err.Error()))
}
