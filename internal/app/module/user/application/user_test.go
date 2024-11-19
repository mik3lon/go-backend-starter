package user_application_test

import (
	"context"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"github.com/stretchr/testify/mock"
)

// Mock dependencies
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*user_domain.User, error) {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(*user_domain.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Save(ctx context.Context, user *user_domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindAll(ctx context.Context, page int, size int) (user_domain.UserList, error) {
	args := m.Called(ctx, page, size)
	if users, ok := args.Get(0).(user_domain.UserList); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

type MockIdTokenValidator struct {
	mock.Mock
}

func (m *MockIdTokenValidator) Validate(ctx context.Context, idToken string) (*user_domain.IdTokenClaims, error) {
	args := m.Called(ctx, idToken)
	if claims, ok := args.Get(0).(*user_domain.IdTokenClaims); ok {
		return claims, args.Error(1)
	}
	return nil, args.Error(1)
}

type MockUserEncoder struct {
	mock.Mock
}

func (m *MockUserEncoder) GenerateToken(user *user_domain.User) (*user_domain.TokenDetails, error) {
	args := m.Called(user)
	return args.Get(0).(*user_domain.TokenDetails), args.Error(1)
}
