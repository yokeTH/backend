package apperror

type ErrorStatus string

var (
	StatusFiberError ErrorStatus = "FIBER_ERROR"

	StatusInternalServerError ErrorStatus = "INTERNAL_SERVER_ERROR"
)
