package user_application_test

import (
	"context"
	"errors"
	"testing"

	user_application "github.com/mik3lon/starter-template/internal/app/module/user/application"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"github.com/mik3lon/starter-template/pkg/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateUserProfilePhotoCommandHandler_Handle(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)
	mockUploader := new(MockImageUploader)

	// Create the handler
	handler := user_application.NewUpdateUserProfilePhotoCommandHandler(mockRepo, mockUploader)

	// Define inputs
	fileInfo := &file.FileInfo{
		Filename:    "photo.jpg",
		ContentType: "image/jpeg",
		Size:        1024,
		Content:     []byte{0xFF, 0xD8, 0xFF}, // Example content
	}

	command := &user_application.UpdateUserProfilePhotoCommand{
		Email: "johndoe@example.com",
		Image: fileInfo,
	}

	ctx := context.Background()

	// Mock repository and uploader responses
	existingUser := &user_domain.User{
		ID:                "123",
		Email:             "johndoe@example.com",
		Username:          "johndoe",
		Name:              "John",
		Surname:           "Doe",
		HashedPassword:    "hashedPassword123",
		Role:              "user",
		ProfilePictureUrl: "https://example.com/oldphoto.jpg",
	}

	uploadedImage := &file.UploadFile{
		Url: "https://cdn.example.com/photo.jpg",
	}

	mockRepo.On("FindByEmail", ctx, command.Email).Return(existingUser, nil)
	mockUploader.On("Upload", ctx, *fileInfo).Return(uploadedImage, nil)
	mockRepo.On("Save", ctx, existingUser).Return(nil)

	// Act
	err := handler.Handle(ctx, command)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "FindByEmail", ctx, command.Email)
	mockUploader.AssertCalled(t, "Upload", ctx, *fileInfo)
	mockRepo.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user *user_domain.User) bool {
		return user.ProfilePictureUrl == uploadedImage.Url
	}))
}

func TestUpdateUserProfilePhotoCommandHandler_Handle_InvalidCommand(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)
	mockUploader := new(MockImageUploader)

	// Create the handler
	handler := user_application.NewUpdateUserProfilePhotoCommandHandler(mockRepo, mockUploader)

	// Act
	err := handler.Handle(context.Background(), nil)

	// Assert
	assert.EqualError(t, err, "invalid command")
}

func TestUpdateUserProfilePhotoCommandHandler_Handle_UserNotFound(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)
	mockUploader := new(MockImageUploader)

	// Create the handler
	handler := user_application.NewUpdateUserProfilePhotoCommandHandler(mockRepo, mockUploader)

	// Define inputs
	fileInfo := &file.FileInfo{
		Filename:    "photo.jpg",
		ContentType: "image/jpeg",
		Size:        1024,
		Content:     []byte{0xFF, 0xD8, 0xFF},
	}

	command := &user_application.UpdateUserProfilePhotoCommand{
		Email: "notfound@example.com",
		Image: fileInfo,
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

func TestUpdateUserProfilePhotoCommandHandler_Handle_UploadError(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)
	mockUploader := new(MockImageUploader)

	// Create the handler
	handler := user_application.NewUpdateUserProfilePhotoCommandHandler(mockRepo, mockUploader)

	// Define inputs
	fileInfo := &file.FileInfo{
		Filename:    "photo.jpg",
		ContentType: "image/jpeg",
		Size:        1024,
		Content:     []byte{0xFF, 0xD8, 0xFF},
	}

	command := &user_application.UpdateUserProfilePhotoCommand{
		Email: "johndoe@example.com",
		Image: fileInfo,
	}

	ctx := context.Background()

	// Mock repository and uploader responses
	existingUser := &user_domain.User{
		ID:                "123",
		Email:             "johndoe@example.com",
		Username:          "johndoe",
		Name:              "John",
		Surname:           "Doe",
		HashedPassword:    "hashedPassword123",
		Role:              "user",
		ProfilePictureUrl: "https://example.com/oldphoto.jpg",
	}

	mockRepo.On("FindByEmail", ctx, command.Email).Return(existingUser, nil)
	mockUploader.On("Upload", ctx, *fileInfo).Return(nil, errors.New("upload error"))

	// Act
	err := handler.Handle(ctx, command)

	// Assert
	assert.EqualError(t, err, "upload error")
	mockRepo.AssertCalled(t, "FindByEmail", ctx, command.Email)
	mockUploader.AssertCalled(t, "Upload", ctx, *fileInfo)
}

func TestUpdateUserProfilePhotoCommandHandler_Handle_SaveError(t *testing.T) {
	// Mock dependencies
	mockRepo := new(MockUserRepository)
	mockUploader := new(MockImageUploader)

	// Create the handler
	handler := user_application.NewUpdateUserProfilePhotoCommandHandler(mockRepo, mockUploader)

	// Define inputs
	fileInfo := &file.FileInfo{
		Filename:    "photo.jpg",
		ContentType: "image/jpeg",
		Size:        1024,
		Content:     []byte{0xFF, 0xD8, 0xFF},
	}

	command := &user_application.UpdateUserProfilePhotoCommand{
		Email: "johndoe@example.com",
		Image: fileInfo,
	}

	ctx := context.Background()

	// Mock repository and uploader responses
	existingUser := &user_domain.User{
		ID:                "123",
		Email:             "johndoe@example.com",
		Username:          "johndoe",
		Name:              "John",
		Surname:           "Doe",
		HashedPassword:    "hashedPassword123",
		Role:              "user",
		ProfilePictureUrl: "https://example.com/oldphoto.jpg",
	}

	uploadedImage := &file.UploadFile{
		Url: "https://cdn.example.com/photo.jpg",
	}

	mockRepo.On("FindByEmail", ctx, command.Email).Return(existingUser, nil)
	mockUploader.On("Upload", ctx, *fileInfo).Return(uploadedImage, nil)
	mockRepo.On("Save", ctx, existingUser).Return(errors.New("save error"))

	// Act
	err := handler.Handle(ctx, command)

	// Assert
	assert.EqualError(t, err, "save error")
	mockRepo.AssertCalled(t, "FindByEmail", ctx, command.Email)
	mockUploader.AssertCalled(t, "Upload", ctx, *fileInfo)
	mockRepo.AssertCalled(t, "Save", ctx, existingUser)
}
