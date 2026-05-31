package client

import (
	"fmt"
	"net/http"

	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
)

func mapStatusError(service string, statusCode int) error {
	switch statusCode {
	case http.StatusBadRequest:
		return types.NewInvalidInputError()
	case http.StatusUnauthorized, http.StatusForbidden:
		return types.NewUnauthorizedError()
	case http.StatusNotFound:
		return types.NewNoTFoundOrNoRecordError()
	default:
		if statusCode >= http.StatusInternalServerError {
			return types.NewInternalServerError()
		}
		return fmt.Errorf("%s returned status %d", service, statusCode)
	}
}
