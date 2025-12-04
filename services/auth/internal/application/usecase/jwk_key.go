package usecase

import "github.com/yokeTH/backend/services/auth/internal/application/port"

type JWKKeyUsecase struct {
	generator port.JWKGenerator
}

func NewJWKKeyUsecase(generator port.JWKGenerator) *JWKKeyUsecase {
	uc := JWKKeyUsecase{
		generator: generator,
	}

	return &uc
}
