package user_application_test

import (
	"context"
	"errors"
	"testing"

	user_application "github.com/mik3lon/starter-template/internal/app/module/user/application"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"github.com/stretchr/testify/assert"
)

func TestUserPasswordSignInQueryHandler_Handle(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)
	mockEncoder := new(MockUserEncoder)
	mockEncrypter := new(MockPasswordEncrypter)

	// Create the handler
	handler := user_application.NewUserPasswordSignInQueryHandler(mockRepo, mockEncoder, mockEncrypter)

	// Define inputs
	query := &user_application.UserPasswordSignInQuery{
		Email:    "johndoe@example.com",
		Password: "password123",
	}

	ctx := context.Background()

	// Mock repository response
	existingUser := &user_domain.User{
		ID:             "123",
		Email:          "johndoe@example.com",
		Username:       "johndoe",
		Name:           "John",
		Surname:        "Doe",
		HashedPassword: "hashedPassword123",
		Role:           "user",
	}

	tokenDetails := &user_domain.TokenDetails{
		UserEmail:           existingUser.Email,
		AccessToken:         "access-token",
		RefreshToken:        "refresh-token",
		AccessTokenExpires:  3600,
		RefreshTokenExpires: 7200,
	}

	mockRepo.On("FindByEmail", ctx, query.Email).Return(existingUser, nil)
	mockEncrypter.On("VerifyPassword", existingUser.HashedPassword, query.Password).Return(nil)
	mockEncoder.On("GenerateToken", existingUser).Return(tokenDetails, nil)

	// Act
	result, err := handler.Handle(ctx, query)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "FindByEmail", ctx, query.Email)
	mockEncrypter.AssertCalled(t, "VerifyPassword", existingUser.HashedPassword, query.Password)
	mockEncoder.AssertCalled(t, "GenerateToken", existingUser)

	// Validate the result
	assert.Equal(t, tokenDetails, result)
}

func TestUserPasswordSignInQueryHandler_Handle_InvalidQuery(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)
	mockEncoder := new(MockUserEncoder)
	mockEncrypter := new(MockPasswordEncrypter)

	// Create the handler
	handler := user_application.NewUserPasswordSignInQueryHandler(mockRepo, mockEncoder, mockEncrypter)

	// Act
	result, err := handler.Handle(context.Background(), nil)

	// Assert
	assert.EqualError(t, err, "invalid query")
	assert.Nil(t, result)
}

func TestUserPasswordSignInQueryHandler_Handle_UserNotFound(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)
	mockEncoder := new(MockUserEncoder)
	mockEncrypter := new(MockPasswordEncrypter)

	// Create the handler
	handler := user_application.NewUserPasswordSignInQueryHandler(mockRepo, mockEncoder, mockEncrypter)

	// Define inputs
	query := &user_application.UserPasswordSignInQuery{
		Email:    "notfound@example.com",
		Password: "password123",
	}

	ctx := context.Background()

	// Mock repository response
	mockRepo.On("FindByEmail", ctx, query.Email).Return(nil, errors.New("user not found"))

	// Act
	result, err := handler.Handle(ctx, query)

	// Assert
	assert.EqualError(t, err, "user not found")
	assert.Nil(t, result)
	mockRepo.AssertCalled(t, "FindByEmail", ctx, query.Email)
}

func TestUserPasswordSignInQueryHandler_Handle_InvalidPassword(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)
	mockEncoder := new(MockUserEncoder)
	mockEncrypter := new(MockPasswordEncrypter)

	// Create the handler
	handler := user_application.NewUserPasswordSignInQueryHandler(mockRepo, mockEncoder, mockEncrypter)

	// Define inputs
	query := &user_application.UserPasswordSignInQuery{
		Email:    "johndoe@example.com",
		Password: "wrongpassword",
	}

	ctx := context.Background()

	// Mock repository and encrypter responses
	existingUser := &user_domain.User{
		ID:             "123",
		Email:          "johndoe@example.com",
		HashedPassword: "hashedPassword123",
	}

	mockRepo.On("FindByEmail", ctx, query.Email).Return(existingUser, nil)
	mockEncrypter.On("VerifyPassword", existingUser.HashedPassword, query.Password).Return(errors.New("invalid password"))

	// Act
	result, err := handler.Handle(ctx, query)

	// Assert
	assert.EqualError(t, err, "invalid password")
	assert.Nil(t, result)
	mockRepo.AssertCalled(t, "FindByEmail", ctx, query.Email)
	mockEncrypter.AssertCalled(t, "VerifyPassword", existingUser.HashedPassword, query.Password)
}

func TestUserPasswordSignInQueryHandler_Handle_TokenGenerationError(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)
	mockEncoder := new(MockUserEncoder)
	mockEncrypter := new(MockPasswordEncrypter)

	// Create the handler
	handler := user_application.NewUserPasswordSignInQueryHandler(mockRepo, mockEncoder, mockEncrypter)

	// Define inputs
	query := &user_application.UserPasswordSignInQuery{
		Email:    "johndoe@example.com",
		Password: "password123",
	}

	ctx := context.Background()

	// Mock repository, encrypter, and encoder responses
	existingUser := &user_domain.User{
		ID:             "123",
		Email:          "johndoe@example.com",
		HashedPassword: "hashedPassword123",
	}

	mockRepo.On("FindByEmail", ctx, query.Email).Return(existingUser, nil)
	mockEncrypter.On("VerifyPassword", existingUser.HashedPassword, query.Password).Return(nil)
	mockEncoder.On("GenerateToken", existingUser).Return(nil, errors.New("token generation error"))

	// Act
	result, err := handler.Handle(ctx, query)

	// Assert
	assert.EqualError(t, err, "token generation error")
	assert.Nil(t, result)
	mockRepo.AssertCalled(t, "FindByEmail", ctx, query.Email)
	mockEncrypter.AssertCalled(t, "VerifyPassword", existingUser.HashedPassword, query.Password)
	mockEncoder.AssertCalled(t, "GenerateToken", existingUser)
}
