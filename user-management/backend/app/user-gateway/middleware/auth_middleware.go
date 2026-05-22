package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"tmossDev.github.com/eco-system/shared-components/backend/package/user/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/utils"
	"tmossDev.github.com/eco-system/user-management/backend/domain/user/service"
)

func getJwtTokenFromSession(context *fiber.Ctx) (string, error) {
	authHeader := context.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], nil
		}
		return "", errors.New("authorization header format must be 'Bearer {token}'")
	}

	cookie := context.Cookies("jwt")
	if cookie != "" {
		return cookie, nil
	}

	return "", errors.New("jwt token not found in headers or cookies")
}

func IsAuthenticated(c *fiber.Ctx, service service.PublicUserService) error {
	jwt := c.Cookies("jwt")

	err := service.IsAuthenticated(jwt)
	if err != nil {
		return err
	}

	jwt, err = getJwtTokenFromSession(c)
	if err != nil {
		return err
	}
	issuer, err := utils.GetIssuerFromJwt(jwt, constants.PASSWORD_SECRET_HASHING_KEY)
	if err != nil {
		return errors.New("could not get issuer from jwt")
	}

	c.Set("userId", issuer)

	return c.Next()
}
