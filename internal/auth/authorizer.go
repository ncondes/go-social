package auth

import "github.com/ncondes/go/social/internal/domain"

type Resource interface {
	GetOwnerID() any
}

type Permission string

const (
	PermissionCreatePost Permission = "post:create"
	PermissionReadPost   Permission = "post:read"
	PermissionUpdatePost Permission = "post:update"
	PermissionDeletePost Permission = "post:delete"

	PermissionCreateComment Permission = "comment:create"
	PermissionUpdateComment Permission = "comment:update"
	PermissionDeleteComment Permission = "comment:delete"
)

type Policy struct {
	Permission   Permission
	MinimumLevel int
	AllowOwner   bool
	CustomCheck  func(user *domain.User, resource Resource) bool
}

type Authorizer struct {
	policies map[Permission]Policy
}

func NewAuthorizer() *Authorizer {
	return &Authorizer{
		policies: map[Permission]Policy{
			PermissionUpdatePost: {
				MinimumLevel: domain.RoleLevelModerator,
				AllowOwner:   true,
			},
			PermissionDeletePost: {
				MinimumLevel: domain.RoleLevelAdmin,
				AllowOwner:   true,
			},

			PermissionUpdateComment: {
				MinimumLevel: domain.RoleLevelModerator,
				AllowOwner:   true,
			},
			PermissionDeleteComment: {
				MinimumLevel: domain.RoleLevelAdmin,
				AllowOwner:   true,
			},
		},
	}
}

func (a *Authorizer) Authorize(user *domain.User, permission Permission, resource Resource) bool {
	policy, exists := a.policies[permission]
	if !exists {
		return false
	}

	// Custom check takes precedence
	if policy.CustomCheck != nil {
		return policy.CustomCheck(user, resource)
	}

	// Check ownership if allowed
	if policy.AllowOwner && resource != nil {
		ownerID := resource.GetOwnerID()

		switch value := ownerID.(type) {
		case int64:
			if value == user.ID {
				return true
			}
		}
	}

	// Check role permission level
	return user.Role.Level >= policy.MinimumLevel
}

func (a *Authorizer) CanUpdatePost(user *domain.User, post *domain.Post) bool {
	return a.Authorize(user, PermissionUpdatePost, post)
}

func (a *Authorizer) CanDeletePost(user *domain.User, post *domain.Post) bool {
	return a.Authorize(user, PermissionDeletePost, post)
}
