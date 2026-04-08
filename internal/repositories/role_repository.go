package repositories

import (
	"context"
	"database/sql"

	"github.com/ncondes/go/social/internal/domain"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `
		SELECT id, name, description, level
		FROM roles
		WHERE name = $1
		LIMIT 1
	`

	var role domain.Role
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.Level,
	)

	if err != nil {
		return nil, handleDBError(err, resourceRole)
	}

	return &role, nil
}

var _ domain.RoleRepositoryInterface = (*RoleRepository)(nil)
