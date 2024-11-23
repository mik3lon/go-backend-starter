package user_application_test

import (
	"context"
	"errors"
	"testing"

	user_application "github.com/mik3lon/starter-template/internal/app/module/user/application"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateUserProfileCommandHandler_Handle(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)

	// Create the handler
	handler := user_application.NewUpdateUserProfileCommandHandler(mockRepo)

	// Define inputs
	command := &user_application.UpdateUserProfileCommand{
		Email:    "johndoe@example.com",
		Username: "newusername",
		Name:     "NewName",
		Surname:  "NewSurname",
	}

	ctx := context.Background()

	// Mock repository response
	existingUser := &user_domain.User{
		ID:                "123",
		Email:             "johndoe@example.com",
		Username:          "oldusername",
		Name:              "OldName",
		Surname:           "OldSurname",
		HashedPassword:    "hashedPassword123",
		Role:              "user",
		ProfilePictureUrl: "https://example.com/john.jpg",
	}

	mockRepo.On("FindByEmail", ctx, command.Email).Return(existingUser, nil)
	mockRepo.On("Save", ctx, existingUser).Return(nil)

	// Act
	err := handler.Handle(ctx, command)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "FindByEmail", ctx, command.Email)
	mockRepo.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user *user_domain.User) bool {
		return user.Username == command.Username &&
			user.Name == command.Name &&
			user.Surname == command.Surname
	}))
}

func TestUpdateUserProfileCommandHandler_Handle_InvalidCommand(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)

	// Create the handler
	handler := user_application.NewUpdateUserProfileCommandHandler(mockRepo)

	// Act
	err := handler.Handle(context.Background(), nil)

	// Assert
	assert.EqualError(t, err, "invalid command")
}

func TestUpdateUserProfileCommandHandler_Handle_UserNotFound(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)

	// Create the handler
	handler := user_application.NewUpdateUserProfileCommandHandler(mockRepo)

	// Define inputs
	command := &user_application.UpdateUserProfileCommand{
		Email:    "notfound@example.com",
		Username: "newusername",
		Name:     "NewName",
		Surname:  "NewSurname",
	}

	ctx := context.Background()

	// Mock repository response
	mockRepo.On("FindByEmail", ctx, command.Email).Return(nil, errors.New("user not found"))

	// Act
	err := handler.Handle(ctx, command)

	// Assert
	assert.EqualError(t, err, "user not found")
	mockRepo.AssertCalled(t, "FindByEmail", ctx, command.Email)
}

func TestUpdateUserProfileCommandHandler_Handle_SaveError(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)

	// Create the handler
	handler := user_application.NewUpdateUserProfileCommandHandler(mockRepo)

	// Define inputs
	command := &user_application.UpdateUserProfileCommand{
		Email:    "johndoe@example.com",
		Username: "newusername",
		Name:     "NewName",
		Surname:  "NewSurname",
	}

	ctx := context.Background()

	// Mock repository response
	existingUser := &user_domain.User{
		ID:                "123",
		Email:             "johndoe@example.com",
		Username:          "oldusername",
		Name:              "OldName",
		Surname:           "OldSurname",
		HashedPassword:    "hashedPassword123",
		Role:              "user",
		ProfilePictureUrl: "https://example.com/john.jpg",
	}

	mockRepo.On("FindByEmail", ctx, command.Email).Return(existingUser, nil)
	mockRepo.On("Save", ctx, existingUser).Return(errors.New("save error"))

	// Act
	err := handler.Handle(ctx, command)

	// Assert
	assert.EqualError(t, err, "save error")
	mockRepo.AssertCalled(t, "FindByEmail", ctx, command.Email)
	mockRepo.AssertCalled(t, "Save", ctx, existingUser)
}
