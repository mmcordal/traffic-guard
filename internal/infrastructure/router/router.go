package router

import (
	"traffic-guarder/internal/infrastructure/app"
	"traffic-guarder/internal/infrastructure/errorsx"

	"github.com/gofiber/fiber/v2"
)

func Get(r fiber.Router, path string, handlers func(ctx *app.Ctx) errorsx.APIError) {
	r.Get(path, func(c *fiber.Ctx) error {
		if err := handlers(&app.Ctx{Ctx: c}); err != nil {
			return err
		}
		return nil
	})
}

func Post(r fiber.Router, path string, handlers func(ctx *app.Ctx) errorsx.APIError) {
	r.Post(path, func(c *fiber.Ctx) error {
		if err := handlers(&app.Ctx{Ctx: c}); err != nil {
			return err
		}
		return nil
	})
}

func Put(r fiber.Router, path string, handlers func(ctx *app.Ctx) errorsx.APIError) {
	r.Put(path, func(c *fiber.Ctx) error {
		if err := handlers(&app.Ctx{Ctx: c}); err != nil {
			return err
		}
		return nil
	})
}

func Delete(r fiber.Router, path string, handlers func(ctx *app.Ctx) errorsx.APIError) {
	r.Delete(path, func(c *fiber.Ctx) error {
		if err := handlers(&app.Ctx{Ctx: c}); err != nil {
			return err
		}
		return nil
	})
}
