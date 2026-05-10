package middleware

import (
	"github.com/gofiber/fiber/v2"
	"tmossDev.github.com/eco-system/user-management/backend/package/user/service"
)

func IsAuthorized(c *fiber.Ctx, page string, service service.PublicUserService) error {
	jwt := c.Cookies("jwt")

	err := service.IsAuthorized(jwt, page)
	if err != nil {
		return err
	}

	return nil
}
