package tracer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type requestIDKey struct{}

func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(requestIDKey{}, uuid.New().String())
		return c.Next()
	}
}

func GetRequestID(ctx *fiber.Ctx) string {
	if ctx == nil {
		return ""
	}
	if reID, ok := ctx.Locals(requestIDKey{}).(string); ok {
		return reID
	}
	return ""
}
