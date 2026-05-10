package types

import "tmossDev.github.com/eco-system/shared-components/backend/package/constants"

type SocketError struct {
	statusCode int
	message    string
}

func NewSocketError(statusCode int, message string) *SocketError {
	return &SocketError{
		statusCode: statusCode,
		message:    message,
	}
}

func (e *SocketError) Error() string {
	return e.message
}

func (e *SocketError) StatusCode() int {
	return e.statusCode
}

func NewInternalServerError() error {
	return NewSocketError(constants.InternalServerErrorCode, constants.InternalServerErrorName)
}

func NewBadRequestError() error {
	return NewSocketError(constants.BadRequestCode, constants.BadRequestName)
}

func NewInvalidInputError() error {
	return NewSocketError(constants.InvalidInputCode, constants.InvalidInputErrorName)
}

func NewNotImplementedError() error {
	return NewSocketError(constants.NotImplementedCode, constants.NotImplementedErrorName)
}

func NewNoTFoundOrNoRecordError() error {
	return NewSocketError(constants.NotFoundCode, constants.NotFoundErrorName)
}

func NewUnauthorizedError() error {
	return NewSocketError(constants.UnauthorizedCode, constants.UnauthorizedRequestName)
}
