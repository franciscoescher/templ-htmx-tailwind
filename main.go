package main

import (
	"platform/internal/components"
	"platform/internal/components/base"
	"platform/internal/logger"
	"platform/internal/tracer"
	"platform/internal/utils"
	"reflect"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"go.uber.org/dig"
)

func main() {
	c := dig.New()
	Provide(
		c,
		logger.NewLogger,
		BuildApp,
		func() []base.MenuItem {
			return []base.MenuItem{
				{Title: "Home", URL: "/", Component: components.Home()},
				{Title: "About", URL: "/about", Component: components.About()},
				{Title: "Contact", URL: "/contact", Component: components.Contact()},
			}
		},
		base.New,
	)

	err := c.Invoke(func(app *fiber.App, l *logger.Logger) error {
		l.Info(nil, "starting application")
		err := app.Listen(":8000")
		if err != nil {
			l.Error(nil, "error listening to port", err)
		}
		return err
	})
	if err != nil {
		panic(err)
	}
}

func Provide(c *dig.Container, constructors ...any) {
	for _, constructor := range constructors {
		opts := make([]dig.ProvideOption, 0)
		if reflect.TypeOf(constructor).Out(0).Implements(reflect.TypeOf((*utils.Controller)(nil)).Elem()) {
			opts = append(opts, dig.Group("controller"))
		}
		err := c.Provide(constructor, opts...)
		if err != nil {
			panic(err)
		}
	}
}

type Params struct {
	dig.In

	Logger      *logger.Logger
	Controllers []utils.Controller `group:"controller"`
}

func BuildApp(p Params) *fiber.App {
	f := fiber.New()
	f.Use(recover.New())
	f.Use(tracer.Middleware())
	f.Use(logger.Middleware(p.Logger))

	// Static files
	f.Static("/src", "./src/dist")

	// Register routes
	for _, ctrl := range p.Controllers {
		ctrl.Register(f)
	}
	// Custom 404 handler
	f.Use(utils.RenderForFiber(func(ctx *fiber.Ctx) (templ.Component, error) {
		return components.NotFound(), nil
	}))

	return f
}

// AddComponent is a helper function that renders a templ component to the fiber response.
func AddComponent(f *fiber.App, path string, component templ.Component) {
	f.Get(path, utils.RenderForFiber(func(ctx *fiber.Ctx) (templ.Component, error) {
		return component, nil
	}))
}
