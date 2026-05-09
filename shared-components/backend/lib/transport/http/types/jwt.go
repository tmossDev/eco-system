package types

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTConfig holds the configuration for JWT middleware
type JWTConfig struct {
	SecretKey     []byte
	TokenExpiry   time.Duration
	SigningMethod jwt.SigningMethod
	TokenPrefix   string
}

// CustomClaims extends jwt.StandardClaims to include user-specific claims
type CustomClaims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role,omitempty"`
	jwt.StandardClaims
}
