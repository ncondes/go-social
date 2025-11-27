package dtos

type CursorBasedPaginationResponseDTO[T any] struct {
	Data       []T                       `json:"data"`
	Pagination CursorBasedPaginationMeta `json:"pagination"`
}

type CursorBasedPaginationMeta struct {
	Limit      int    `json:"limit"`
	NextCursor string `json:"next_cursor"`
}
