package types_test

import (
	"testing"

	"tmossDev.github.com/eco-system/shared-components/backend/package/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
)

func commonTestSocketErrorFlow(t *testing.T, err error, expectedStatusCode int, expectedMessage string) {
	if socketErr, ok := err.(*types.SocketError); ok {
		statusCode := socketErr.StatusCode()
		if expectedStatusCode != statusCode {
			t.Errorf("Expected Status Code to be '%d' but got '%d'", expectedStatusCode, statusCode)
		}

		errorDescription := socketErr.Error()
		if expectedMessage != errorDescription {
			t.Errorf("Expected Message to be '%s' but got '%s'", expectedMessage, errorDescription)
		}
	} else {
		t.Error("Expected Socket Error Type")
	}
}

func TestNewSocketError_withNewInternalServerError_expectConstantsToMatch(t *testing.T) {
	err := types.NewInternalServerError()
	commonTestSocketErrorFlow(t, err, constants.InternalServerErrorCode, constants.InternalServerErrorName)
}

func TestNewSocketError_withNewBadRequestError_expectConstantsToMatch(t *testing.T) {
	err := types.NewBadRequestError()
	commonTestSocketErrorFlow(t, err, constants.BadRequestCode, constants.BadRequestName)
}

func TestNewSocketError_withTestNewInvalidInputError_expectConstantsToMatch(t *testing.T) {
	err := types.NewInvalidInputError()
	commonTestSocketErrorFlow(t, err, constants.InvalidInputCode, constants.InvalidInputErrorName)
}

func TestNewSocketError_withTestNewNotImplementedError_expectConstantsToMatch(t *testing.T) {
	err := types.NewNotImplementedError()
	commonTestSocketErrorFlow(t, err, constants.NotImplementedCode, constants.NotImplementedErrorName)
}

func TestNewSocketError_withTestNewNoTFoundOrNoRecordError_expectConstantsToMatch(t *testing.T) {
	err := types.NewNoTFoundOrNoRecordError()
	commonTestSocketErrorFlow(t, err, constants.NotFoundCode, constants.NotFoundErrorName)
}
