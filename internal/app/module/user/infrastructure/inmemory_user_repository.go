package user_infrastructure

import (
	"context"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
)

// InMemoryUserRepository is an in-memory implementation of UserRepository.
type InMemoryUserRepository struct {
	usersByEmail map[string]*user_domain.User
}

// NewMemoryUserRepository initializes a new in-memory repository.
func NewMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{usersByEmail: make(map[string]*user_domain.User)}
}

// Save stores a user in memory.
func (r *InMemoryUserRepository) Save(ctx context.Context, user *user_domain.User) error {
	r.usersByEmail[user.Email] = user
	return nil
}

func (r *InMemoryUserRepository) FindByEmail(ctx context.Context, email string) (*user_domain.User, error) {
	user, exists := r.usersByEmail[email]
	if !exists {
		return nil, user_domain.ErrUserNotFound
	}
	return user, nil
}
