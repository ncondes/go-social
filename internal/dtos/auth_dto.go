package dtos

import "github.com/ncondes/go/social/internal/domain"

type ActivateUserDTO struct {
	Token string `json:"token" validate:"required" example:"2fa66d86-ec2b-4766-ae19-bdcd24045d82"`
}

type GenerateTokenDTO struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
}

type GenerateTokenResponseDTO struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
}

type RegisterResponseDTO struct {
	User  *domain.User `json:"user"`
	Token string       `json:"token" example:"2fa66d86-ec2b-4766-ae19-bdcd24045d82"`
}
