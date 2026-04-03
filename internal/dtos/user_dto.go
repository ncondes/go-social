package dtos

type AuthorInfoDTO struct {
	ID       int64  `json:"id"       example:"1"`
	Fullname string `json:"fullname" example:"John Doe"`
	Username string `json:"username" example:"johndoe"`
}
