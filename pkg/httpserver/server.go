package httpserver

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/yokeTH/backend/pkg/apperror"
)

type Server struct {
	config *Config
	*fiber.App
}

func New(opts ...ServerOption) *Server {

	cfg := defaultConfig

	server := &Server{
		config: &cfg,
	}

	for _, opt := range opts {
		opt(server)
	}

	app := fiber.New(fiber.Config{
		AppName:       server.config.Name,
		BodyLimit:     server.config.BodyLimitMB * 1024 * 1024,
		CaseSensitive: true,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		ErrorHandler:  apperror.ErrorHandler,
	})

	app.Use(requestid.New())

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	server.App = app

	return server
}
