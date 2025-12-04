package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/rs/zerolog/log"
	"github.com/yokeTH/backend/pkg/apperror"
	"github.com/yokeTH/backend/pkg/resp"
)

type serverInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Env     string `json:"env"`
}

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

func (s *Server) Start(ctx context.Context, stop context.CancelFunc) {
	version := os.Getenv("APP_VERSION")
	if version == "" {
		version = "unknown"
	}

	s.App.Get("/", func(ctx fiber.Ctx) error {
		return ctx.JSON(resp.Success(serverInfo{
			Name:    s.config.Name,
			Version: version,
			Env:     s.config.Env,
		}))
	})

	go func() {
		if err := s.Listen(fmt.Sprintf(":%d", s.config.Port)); err != nil {
			log.Error().Err(err).Msg("failed to start server")
			stop()
		}
	}()

	defer func() {
		if err := s.Shutdown(); err != nil {
			log.Error().Err(err).Msg("failed to shutdown server")
		}
	}()

	<-ctx.Done()

	log.Info().Msg("shutting down server...")
}
