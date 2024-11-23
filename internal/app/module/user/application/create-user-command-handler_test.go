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

func TestCreateUserCommandHandler_Handle(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)
	mockEncrypter := new(MockPasswordEncrypter)

	// Create the handler
	handler := user_application.NewCreateUserCommandHandler(mockRepo, mockEncrypter)

	// Define inputs
	command := &user_application.CreateUserCommand{
		ID:                "123",
		Name:              "John",
		Surname:           "Doe",
		Username:          "johndoe",
		PlainPassword:     "password123",
		Email:             "johndoe@example.com",
		Role:              "user",
		ProfilePictureUrl: "https://example.com/john.jpg",
		IsFormSocialAuth:  false,
	}

	ctx := context.Background()

	// Expected outputs
	hashedPassword := "hashedPassword123"
	mockEncrypter.On("GenerateHashedPassword", command.IsFormSocialAuth, command.PlainPassword).
		Return(hashedPassword, nil)

	mockRepo.On("Save", ctx, mock.AnythingOfType("*user_domain.User")).Return(nil)

	// Act
	err := handler.Handle(ctx, command)

	// Assert
	assert.NoError(t, err)
	mockEncrypter.AssertCalled(t, "GenerateHashedPassword", command.IsFormSocialAuth, command.PlainPassword)
	mockRepo.AssertCalled(t, "Save", ctx, mock.AnythingOfType("*user_domain.User"))

	// Verify the created user object
	mockRepo.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user *user_domain.User) bool {
		return user.ID == command.ID &&
			user.Name == command.Name &&
			user.Surname == command.Surname &&
			user.Username == command.Username &&
			user.Email == command.Email &&
			user.HashedPassword == hashedPassword &&
			user.Role == command.Role &&
			user.ProfilePictureUrl == command.ProfilePictureUrl
	}))
}

func TestCreateUserCommandHandler_Handle_InvalidCommand(t *testing.T) {
	// Create handler
	handler := user_application.NewCreateUserCommandHandler(nil, nil)

	// Act
	err := handler.Handle(context.Background(), nil)

	// Assert
	assert.EqualError(t, err, "invalid command")
}

func TestCreateUserCommandHandler_Handle_HashingError(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)
	mockEncrypter := new(MockPasswordEncrypter)

	// Create the handler
	handler := user_application.NewCreateUserCommandHandler(mockRepo, mockEncrypter)

	// Define inputs
	command := &user_application.CreateUserCommand{
		ID:                "123",
		Name:              "John",
		Surname:           "Doe",
		Username:          "johndoe",
		PlainPassword:     "password123",
		Email:             "johndoe@example.com",
		Role:              "user",
		ProfilePictureUrl: "https://example.com/john.jpg",
		IsFormSocialAuth:  false,
	}

	ctx := context.Background()

	// Expected outputs
	mockEncrypter.On("GenerateHashedPassword", command.IsFormSocialAuth, command.PlainPassword).
		Return("", errors.New("hashing error"))

	// Act
	err := handler.Handle(ctx, command)

	// Assert
	assert.EqualError(t, err, "failed to generate hashed password")
	mockEncrypter.AssertCalled(t, "GenerateHashedPassword", command.IsFormSocialAuth, command.PlainPassword)
	mockRepo.AssertNotCalled(t, "Save", ctx, mock.Anything)
}

func TestCreateUserCommandHandler_Handle_SaveError(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)
	mockEncrypter := new(MockPasswordEncrypter)

	// Create the handler
	handler := user_application.NewCreateUserCommandHandler(mockRepo, mockEncrypter)

	// Define inputs
	command := &user_application.CreateUserCommand{
		ID:                "123",
		Name:              "John",
		Surname:           "Doe",
		Username:          "johndoe",
		PlainPassword:     "password123",
		Email:             "johndoe@example.com",
		Role:              "user",
		ProfilePictureUrl: "https://example.com/john.jpg",
		IsFormSocialAuth:  false,
	}

	ctx := context.Background()

	// Expected outputs
	hashedPassword := "hashedPassword123"
	mockEncrypter.On("GenerateHashedPassword", command.IsFormSocialAuth, command.PlainPassword).
		Return(hashedPassword, nil)

	mockRepo.On("Save", ctx, mock.AnythingOfType("*user_domain.User")).Return(errors.New("save error"))

	// Act
	err := handler.Handle(ctx, command)

	// Assert
	assert.EqualError(t, err, "save error")
	mockEncrypter.AssertCalled(t, "GenerateHashedPassword", command.IsFormSocialAuth, command.PlainPassword)
	mockRepo.AssertCalled(t, "Save", ctx, mock.AnythingOfType("*user_domain.User"))
}
