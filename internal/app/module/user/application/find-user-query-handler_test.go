package user_application_test

import (
	"context"
	"errors"
	"testing"

	user_application "github.com/mik3lon/starter-template/internal/app/module/user/application"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"github.com/stretchr/testify/assert"
)

func TestFindUserQueryHandler_Handle(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)

	// Create the handler
	handler := user_application.NewFindUserQueryHandler(mockRepo)

	// Define inputs
	query := &user_application.FindUserQuery{
		Email: "johndoe@example.com",
	}

	ctx := context.Background()

	// Mock repository response
	expectedUser := &user_domain.User{
		ID:                "123",
		Email:             "johndoe@example.com",
		Name:              "John",
		Surname:           "Doe",
		Username:          "johndoe",
		HashedPassword:    "hashedPassword123",
		Role:              "user",
		ProfilePictureUrl: "https://example.com/john.jpg",
	}

	mockRepo.On("FindByEmail", ctx, query.Email).Return(expectedUser, nil)

	// Act
	result, err := handler.Handle(ctx, query)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "FindByEmail", ctx, query.Email)
	assert.NotNil(t, result)

	// Validate the result structure
	response, ok := result.(*user_application.FindUserResponse)
	assert.True(t, ok)
	assert.Equal(t, expectedUser.ID, response.ID)
	assert.Equal(t, expectedUser.Email, response.Email)
	assert.Equal(t, expectedUser.Name, response.Name)
	assert.Equal(t, expectedUser.Surname, response.Surname)
	assert.Equal(t, expectedUser.Username, response.Username)
	assert.Equal(t, expectedUser.Role, response.Role)
	assert.Equal(t, expectedUser.ProfilePictureUrl, response.ProfilePictureUrl)
}

func TestFindUserQueryHandler_Handle_InvalidQuery(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)

	// Create the handler
	handler := user_application.NewFindUserQueryHandler(mockRepo)

	// Act
	result, err := handler.Handle(context.Background(), nil)

	// Assert
	assert.EqualError(t, err, "invalid query")
	assert.Nil(t, result)
}

func TestFindUserQueryHandler_Handle_UserNotFound(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)

	// Create the handler
	handler := user_application.NewFindUserQueryHandler(mockRepo)

	// Define inputs
	query := &user_application.FindUserQuery{
		Email: "notfound@example.com",
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

func TestFindUserQueryHandler_Handle_RepoError(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)

	// Create the handler
	handler := user_application.NewFindUserQueryHandler(mockRepo)

	// Define inputs
	query := &user_application.FindUserQuery{
		Email: "johndoe@example.com",
	}

	ctx := context.Background()

	// Mock repository response
	mockRepo.On("FindByEmail", ctx, query.Email).Return(nil, errors.New("database error"))

	// Act
	result, err := handler.Handle(ctx, query)

	// Assert
	assert.EqualError(t, err, "database error")
	assert.Nil(t, result)
	mockRepo.AssertCalled(t, "FindByEmail", ctx, query.Email)
}
