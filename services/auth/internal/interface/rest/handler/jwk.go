package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/guregu/dynamo/v2"
	"github.com/yokeTH/backend/pkg/apperror"
	"github.com/yokeTH/backend/services/auth/internal/application/port"
	"github.com/yokeTH/backend/services/auth/internal/domain/repository"
)

type jwkHandler struct {
	keygenerator  port.JWKGenerator
	keyRepository repository.JWKKeyRepository
}

func NewJWKHandler(keygenerator port.JWKGenerator, keyRepository repository.JWKKeyRepository) *jwkHandler {
	return &jwkHandler{
		keygenerator:  keygenerator,
		keyRepository: keyRepository,
	}
}

func (h *jwkHandler) GetPublicJWKS(ctx fiber.Ctx) error {
	pubs, err := h.keyRepository.GetPublicKeys(ctx.Context())
	if err != nil {
		if errors.Is(err, dynamo.ErrNotFound) {
			return apperror.NotFoundError(err, "public jwk keys not found", ErrStatusPublicJWKSNotFound)
		}

		return apperror.InternalServerError(err, "failed to get jwk keys", ErrStatusPublicJWKInternalError)
	}

	ctx.Set("Content-Type", "application/json")
	return ctx.SendString(fmt.Sprintf("{\"jwks\": %s}", pubs))
}
