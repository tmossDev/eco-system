package constants

import "net/http"

const (
	SuccessCode             = http.StatusOK
	SuccessName             = "OK"
	BadRequestCode          = http.StatusBadRequest
	BadRequestName          = "Bad Request"
	UnauthorizedCode        = http.StatusUnauthorized
	UnauthorizedRequestName = "Permission Denied"
	ForbiddenCode           = http.StatusForbidden
	ForbiddenErrorName      = "Access Forbidden"
	NotFoundCode            = http.StatusNotFound
	NotFoundErrorName       = "No Record Found"
	InvalidInputCode        = http.StatusNotAcceptable
	InvalidInputErrorName   = "Invalid Input"
	InternalServerErrorCode = http.StatusInternalServerError
	InternalServerErrorName = "Internal Server Error"
	NotImplementedCode      = http.StatusNotImplemented
	NotImplementedErrorName = "Not Implemented"
)

const (
	PUT    = "PUT"
	POST   = "POST"
	GET    = "GET"
	DELETE = "DELETE"
)
