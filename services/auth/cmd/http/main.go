package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/yokeTH/backend/pkg/dynamoclient"
	"github.com/yokeTH/backend/pkg/httpserver"
	"github.com/yokeTH/backend/services/auth/internal/domain/model"
	"github.com/yokeTH/backend/services/auth/internal/infra/dynamodb"
	"github.com/yokeTH/backend/services/auth/internal/infra/jwk"
	"github.com/yokeTH/backend/services/auth/internal/interface/rest/handler"
)

type appConfig struct {
	Server   httpserver.Config   `envPrefix:"HTTP_SERVER_"`
	Dynamodb dynamoclient.Config `envPrefix:"DYNAMO_"`

	KEKHEX string `env:"KEK_HEX,required"`
}

func newConfigFromEnv() *appConfig {
	config := &appConfig{}

	if err := env.Parse(config); err != nil {
		log.Fatal().Err(err).Msg("Unable to parse env vars: %s")
	}

	return config
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("Unable to load .env file: %s")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	cfg := newConfigFromEnv()

	serverCfg := cfg.Server

	s := httpserver.New(
		httpserver.WithConfig(&serverCfg),
	)

	localEncryptor, err := jwk.NewLocalKeyEncryptorFromHex(cfg.KEKHEX)
	if err != nil {
		panic(err)
	}

	db := dynamoclient.NewDynamoClient(ctx, cfg.Dynamodb)
	jwkRepository := dynamodb.NewJWKKeyDynamodb(db)
	jwkCnt, err := jwkRepository.Count(ctx)
	if err != nil {
		panic(err)
	}

	jwkGenerator := jwk.NewJWKGenerator(localEncryptor, "env://KEK_HEX", model.JWKKeyStatusActive, 24*time.Hour*7)

	if jwkCnt == 0 {
		newKey, err := jwkGenerator.Generate(ctx)
		if err != nil {
			panic(err)
		}
		if err := jwkRepository.CreateKey(ctx, newKey); err != nil {
			panic(err)
		}
	}

	jwkHandler := handler.NewJWKHandler(jwkGenerator, jwkRepository)

	wk := s.Group("/.well-known")
	wk.Get("/jwks.json", jwkHandler.GetPublicJWKS)

	s.Start(ctx, stop)
}
