package http

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
	"tmossDev.github.com/eco-system/shared-components/backend/package/transport/http/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/transport/http/types"
)

// NewJWTMiddleware creates a new JWT middleware with custom configuration
func NewJWTMiddleware(config types.JWTConfig) func([]string) iris.Handler {

	return func(escapedRoutes []string) iris.Handler {
		return func(ctx iris.Context) {
			// Check if the current route is in escaped routes
			currentPath := strings.Replace(ctx.Path(), constants.ApiPrefix, "", 1)

			for _, route := range escapedRoutes {
				if isEscapedRoute(currentPath, route) {
					ctx.Next()
					return
				}
			}

			requestID := ctx.Values().GetString(constants.CTXRequestIdKey)

			// Extract token from Authorization header
			tokenString := extractToken(ctx, config)
			if tokenString == "" {
				logger.Infof(requestID, "JWT token is missing")
				RespondWithError(ctx.ResponseWriter(), requestID, constants.ErrUnauthorized)
				return
			}

			// Parse and validate token
			claims, err := validateToken(tokenString, &config)
			if err != nil {
				logger.Infof(requestID, "Invalid JWT token: %v", err)
				RespondWithError(ctx.ResponseWriter(), requestID, constants.ErrUnauthorized)
				return
			}

			// Store claims in context
			ctx.Values().Set("claims", claims)
			ctx.Next()
		}
	}
}

func isEscapedRoute(currentPath string, route string) bool {
	if strings.HasSuffix(route, "/*") {
		prefix := strings.TrimSuffix(route, "*")

		return strings.HasPrefix(strings.ToLower(currentPath), strings.ToLower(prefix))
	}

	return strings.EqualFold(currentPath, route)
}

// Helper functions
func extractToken(ctx iris.Context, config types.JWTConfig) string {
	bearerToken := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(bearerToken, config.TokenPrefix) {
		return ""
	}
	return strings.TrimPrefix(bearerToken, config.TokenPrefix)
}

func validateToken(tokenString string, config *types.JWTConfig) (*types.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &types.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if token.Method != config.SigningMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.SecretKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch {
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				return nil, constants.ErrExpiredToken
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				return nil, constants.ErrMalformedToken
			default:
				return nil, constants.ErrInvalidToken
			}
		}
		return nil, err
	}

	claims, ok := token.Claims.(*types.CustomClaims)
	if !ok || !token.Valid {
		return nil, constants.ErrInvalidClaims
	}

	return claims, nil
}
