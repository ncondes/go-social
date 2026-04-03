package dtos

type DataResponseDTO struct {
	Data any `json:"data"`
}

type ErrorResponseDTO struct {
	Error string `json:"error" example:"something went wrong"`
}

type ErrorsResponseDTO struct {
	Errors []string `json:"errors" example:"field is required"`
}

type CursorBasedPaginationResponseDTO[T any] struct {
	Data       []T                          `json:"data"`
	Pagination CursorBasedPaginationMetaDTO `json:"pagination"`
}

type CursorBasedPaginationMetaDTO struct {
	Limit      int    `json:"limit"`
	NextCursor string `json:"next_cursor"`
}
