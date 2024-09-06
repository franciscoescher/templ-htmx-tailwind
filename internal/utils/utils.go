package utils

import (
	"bytes"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	Register(f *fiber.App)
}

// Define a type constraint interface for allowed function types
type RenderableFunc interface {
	func(ctx *fiber.Ctx) templ.Component | func(ctx *fiber.Ctx) (templ.Component, error)
}

// RenderForFiber is a helper that accepts a function that returns a templ component and renders it to the response.
func RenderForFiber[T RenderableFunc](fn T) func(ctx *fiber.Ctx) error {
	switch f := any(fn).(type) {

	// Handle the case where the function returns (templ.Component, error)
	case func(ctx *fiber.Ctx) (templ.Component, error):
		return renderForFiber(f)

	// Handle the case where the function returns templ.Component
	case func(ctx *fiber.Ctx) templ.Component:
		// Wrap the func(ctx *fiber.Ctx) templ.Component so that it matches the func(ctx *fiber.Ctx) (templ.Component, error) signature
		f2 := func(ctx *fiber.Ctx) (templ.Component, error) {
			return f(ctx), nil
		}
		return renderForFiber(f2)

	default:
		// This should never happen if RenderableFunc is properly constrained
		return func(ctx *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusInternalServerError, "Unsupported function type")
		}
	}
}

// renderForFiber handles the actual rendering of a templ.Component to the Fiber response.
func renderForFiber(f func(ctx *fiber.Ctx) (templ.Component, error)) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		// Set the Content-Type to HTML
		ctx.Response().Header.Set("Content-Type", "text/html")

		// Call the render function to get the component and error
		component, err := f(ctx)
		if err != nil {
			return err
		}

		// Create a buffer to render the component's HTML output
		var buf bytes.Buffer

		// Render the component into the buffer
		err = component.Render(ctx.Context(), &buf)
		if err != nil {
			return err
		}

		// Write the rendered HTML content to the response body
		ctx.Response().SetBody(buf.Bytes())

		return nil
	}
}
