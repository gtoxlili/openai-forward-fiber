package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"strings"
)

func Error(c *fiber.Ctx) error {
	err := c.Next()
	if err != nil {
		log.Warn().Err(err).Msg("Server Error")
		// 根本错误
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			unwrapped = err
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"message": unwrapped.Error(),
				"code":    strings.Replace(err.Error(), ": "+unwrapped.Error(), "", 1),
			},
		})
	}
	return nil
}
