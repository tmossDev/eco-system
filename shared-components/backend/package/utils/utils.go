package utils

import (
	"tmossDev.github.com/eco-system/shared-components/backend/package/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
)

func requestIdOrDefault(requestIds ...string) string {
	if len(requestIds) > 0 && requestIds[0] != "" {
		return requestIds[0]
	}

	return constants.DefaultRequestId
}

func LogPreparingError(queryName string, err error, requestIds ...string) {
	logger.Errorf(requestIdOrDefault(requestIds...), "Error preparing '%s' query: %s", queryName, err.Error())
}

func LogExecutingError(queryName string, err error, requestIds ...string) {
	logger.Errorf(requestIdOrDefault(requestIds...), "Error executing '%s' query: %s", queryName, err.Error())
}

func LogBeginingTnxError(queryName string, err error, requestIds ...string) {
	logger.Errorf(requestIdOrDefault(requestIds...), "Error begining transaction '%s' query: %s", queryName, err.Error())
}

func LogCommitError(queryName string, err error, requestIds ...string) {
	logger.Errorf(requestIdOrDefault(requestIds...), "Error commiting transaction '%s' query: %s", queryName, err.Error())
}
