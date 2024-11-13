package user_domain

import "context"

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindAll(ctx context.Context, page int, size int) (UserList, error)
}
