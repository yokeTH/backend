package handler

import "github.com/yokeTH/backend/pkg/apperror"

var (
	ErrStatusPublicJWKSNotFound     apperror.ErrorStatus = "PUBLIC_JWKS_NOT_FOUND"
	ErrStatusPublicJWKInternalError apperror.ErrorStatus = "PUBLIC_JWKS_INTERNAL_ERROR"
)
