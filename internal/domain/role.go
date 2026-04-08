package domain

import (
	"context"
	"errors"
)

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

const (
	RoleLevelUser      = 1
	RoleLevelModerator = 2
	RoleLevelAdmin     = 3
)

type RoleRepositoryInterface interface {
	GetByName(ctx context.Context, name string) (*Role, error)
}

var (
	ErrRoleNotFound = errors.New("role not found")
)
