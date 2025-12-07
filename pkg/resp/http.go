package resp

type SuccessResponse[T any] struct {
	Data       T          `json:"data"`
	Pagination Pagination `json:"pagination,omitzero"`
}

type Pagination struct {
	CurrentPage int `json:"current" example:"1"`
	LastPage    int `json:"last" example:"4"`
	Limit       int `json:"limit" exaimple:"10"`
	Total       int `json:"total" example:"36"`
}

func Success[T any](data T) *SuccessResponse[T] {
	return &SuccessResponse[T]{Data: data}
}

func (s *SuccessResponse[T]) WithPagination(currentPage int, lastPage int, limit int, total int) *SuccessResponse[T] {
	s.Pagination = Pagination{
		CurrentPage: currentPage,
		LastPage:    lastPage,
		Limit:       limit,
		Total:       total,
	}

	return s
}
