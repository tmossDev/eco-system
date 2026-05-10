package constants

import (
	"errors"
	"net/http"
)

// General errors
var (
	ErrBadPageSize  = errors.New("invalid page size")
	ErrBadPageIndex = errors.New("invalid page index")
	ErrUnauthorized = errors.New("unauthorized access")

	ErrInternalExceptionCode = "ERR_INTERNAL_EXCEPTION"
	ErrInternalExceptionDesc = "An internal server error occurred."
)

var BadPayload = "ERR_BAD_PAYLOAD_FIELDS"
var SuccessCode = "00"
var SuccessMessage = "Request performed successfully"

// JWT errors
var ErrInvalidToken = errors.New("invalid token")
var ErrExpiredToken = errors.New("token has expired")
var ErrMalformedToken = errors.New("malformed token")
var ErrMissingToken = errors.New("missing token")
var ErrInvalidClaims = errors.New("invalid token claims")

var (
	ErrCodeMap = map[error]string{
		ErrInvalidToken:   "ERR_BAD_TOKEN",
		ErrExpiredToken:   "ERR_EXPIRED_TOKEN",
		ErrMalformedToken: "ERR_MALFORMED_TOKEN",
		ErrMissingToken:   "ERR_MISSING_TOKEN",
		ErrInvalidClaims:  "ERR_BAD_TOKEN_CLAIMS",
		ErrBadPageSize:    "ERR_BAD_PAGE_SIZE",
		ErrBadPageIndex:   "ERR_BAD_PAGE_INDEX",
		ErrUnauthorized:   "ERR_UNAUTHORIZED",
	}

	ErrDescriptionMap = map[error]string{
		ErrUnauthorized:   "Unauthorized access.",
		ErrInvalidToken:   "The token provided is invalid.",
		ErrExpiredToken:   "The token provided has expired.",
		ErrMalformedToken: "The token provided is malformed.",
		ErrMissingToken:   "No token provided.",
		ErrInvalidClaims:  "The token claims are invalid.",
		ErrBadPageSize:    "The page size provided is invalid.",
		ErrBadPageIndex:   "The page index provided is invalid.",
	}

	ErrToHTTPStatus = map[error]int{
		ErrUnauthorized: http.StatusUnauthorized,
		ErrBadPageSize:  http.StatusBadRequest,
		ErrBadPageIndex: http.StatusBadRequest,
	}
)
