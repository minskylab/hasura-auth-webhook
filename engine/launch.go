package engine

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (e *Engine) PublicFiberApp() *fiber.App {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: true,
		ExposeHeaders:    "",
		MaxAge:           0,
	}))

	app.Post("/login", e.PublicService.Login)
	app.Post("/register", e.PublicService.Register)
	app.Post("/refresh", e.PublicService.RefreshToken)
	app.Get("/recover-password", e.PublicService.RecoverPassword)
	app.Post("/change-password", e.PublicService.ChangePassword)

	return app
}

func (e *Engine) InternalFiberApp() *fiber.App {
	app := fiber.New()
	app.Use(logger.New())
	// app.Use(cors.New(cors.Config{
	// 	Next:             nil,
	// 	AllowOrigins:     "*",
	// 	AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
	// 	AllowHeaders:     "",
	// 	AllowCredentials: true,
	// 	ExposeHeaders:    "",
	// 	MaxAge:           0,
	// }))

	app.Get("/validate", e.InternalService.HasuraWebhook)
	// app.Get("/users", e.InternalService.ListUsers)
	app.Get("/me", e.InternalService.Me)

	return app
}

func (e *Engine) launchCorrutines(publicErrorSignal, internalErrorSignal chan error) {
	internalServer := e.InternalFiberApp()
	publicServer := e.PublicFiberApp()

	go func(publicErrorSignal chan error) {
		addr := fmt.Sprintf("%s:%d", e.InternalService.Hostname(), e.InternalService.Port())

		if err := internalServer.Listen(addr); err != nil {
			publicErrorSignal <- err
		}
	}(publicErrorSignal)

	go func(internalErrorSignal chan error) {
		addr := fmt.Sprintf("%s:%d", e.PublicService.Hostname(), e.PublicService.Port())

		if err := publicServer.Listen(addr); err != nil {
			publicErrorSignal <- err
		}
	}(internalErrorSignal)
}

func (e *Engine) Launch() chan error {
	combinatedErrorSignal := make(chan error)

	publicErrorSignal := make(chan error)
	internalErrorSignal := make(chan error)

	e.launchCorrutines(publicErrorSignal, internalErrorSignal)

	go func(combinatedErrorSignal, publicErrorSignal, internalErrorSignal chan error) {
		for {
			select {
			case err := <-publicErrorSignal:
				log.WithField("service", "public").Error(err)
				combinatedErrorSignal <- errors.WithMessage(err, "public service error")
			case err := <-internalErrorSignal:
				log.WithField("service", "internal").Error(err)
				combinatedErrorSignal <- errors.WithMessage(err, "internal service error")
			}
		}
	}(combinatedErrorSignal, publicErrorSignal, internalErrorSignal)

	return combinatedErrorSignal
}
