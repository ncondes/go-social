package dtos

type ActivateUserDTO struct {
	Token string `json:"token" validate:"required" example:"2fa66d86-ec2b-4766-ae19-bdcd24045d82"`
}
