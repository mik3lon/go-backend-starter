package user_application_test

import (
	"context"
	"errors"
	user_application "github.com/mik3lon/starter-template/internal/app/module/user/application"
	"testing"

	"github.com/google/uuid"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGoogleSignInQueryHandler_UserFoundSuccessfully(t *testing.T) {
	ctx := context.Background()

	// Arrange
	mockRepo := new(MockUserRepository)
	mockValidator := new(MockIdTokenValidator)
	mockEncoder := new(MockUserEncoder)
	passwordEncrypter := new(MockPasswordEncrypter)

	handler := user_application.NewGoogleSignInQueryHandler(mockRepo, mockValidator, mockEncoder, passwordEncrypter)

	idToken := "test-id-token"
	email := "test@example.com"
	user := &user_domain.User{
		ID:    uuid.NewString(),
		Email: email,
	}
	claims := &user_domain.IdTokenClaims{
		Email: email,
	}
	expectedToken := &user_domain.TokenDetails{
		UserEmail:           email,
		AccessToken:         "mock-access-token",
		RefreshToken:        "mock-refresh-token",
		AccessTokenExpires:  3600,
		RefreshTokenExpires: 7200,
	}

	mockValidator.On("Validate", ctx, idToken).Return(claims, nil)
	mockRepo.On("FindByEmail", ctx, email).Return(user, nil)
	passwordEncrypter.On("GenerateHashedPassword", true, "").Return("encryptedPassword", nil)
	mockEncoder.On("GenerateToken", user).Return(expectedToken, nil)

	// Act
	result, err := handler.Handle(ctx, &user_application.GoogleSignInQuery{IdToken: idToken})

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedToken, result)
	mockValidator.AssertCalled(t, "Validate", ctx, idToken)
	mockRepo.AssertCalled(t, "FindByEmail", ctx, email)
	mockEncoder.AssertCalled(t, "GenerateToken", user)
}

func TestGoogleSignInQueryHandler_UserNotFound_CreatesNewUser(t *testing.T) {
	ctx := context.Background()

	// Arrange
	mockRepo := new(MockUserRepository)
	mockValidator := new(MockIdTokenValidator)
	mockEncoder := new(MockUserEncoder)
	passwordEncrypter := new(MockPasswordEncrypter)

	handler := user_application.NewGoogleSignInQueryHandler(mockRepo, mockValidator, mockEncoder, passwordEncrypter)

	idToken := "test-id-token"
	email := "test@example.com"
	username := "testuser"
	name := "Test"
	surname := "User"
	profilePictureUrl := "https://example.com/profile.jpg"
	claims := &user_domain.IdTokenClaims{
		Email:             email,
		Username:          username,
		Name:              name,
		Surname:           surname,
		ProfilePictureUrl: profilePictureUrl,
	}
	userNotFoundErr := &user_domain.UserNotFound{}
	expectedToken := &user_domain.TokenDetails{
		UserEmail:           email,
		AccessToken:         "mock-access-token",
		RefreshToken:        "mock-refresh-token",
		AccessTokenExpires:  3600,
		RefreshTokenExpires: 7200,
	}

	mockValidator.On("Validate", ctx, idToken).Return(claims, nil)
	passwordEncrypter.On("GenerateHashedPassword", true, "").Return("encryptedPassword", nil)
	mockRepo.On("FindByEmail", ctx, email).Return(nil, userNotFoundErr)
	mockRepo.On("Save", ctx, mock.MatchedBy(func(arg interface{}) bool {
		user, ok := arg.(*user_domain.User)
		if !ok {
			t.Fatalf("Save argument is not a *user_domain.User: %v", arg)
			return false
		}
		return user.Email == email && user.Username == username &&
			user.Name == name && user.Surname == surname &&
			user.ProfilePictureUrl == profilePictureUrl
	})).Return(nil)
	mockEncoder.On("GenerateToken", mock.Anything).Return(expectedToken, nil)

	// Act
	result, err := handler.Handle(ctx, &user_application.GoogleSignInQuery{IdToken: idToken})

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedToken, result)
	mockValidator.AssertCalled(t, "Validate", ctx, idToken)
	mockRepo.AssertCalled(t, "FindByEmail", ctx, email)
	mockEncoder.AssertCalled(t, "GenerateToken", mock.Anything)
}

func TestGoogleSignInQueryHandler_InvalidQuery(t *testing.T) {
	ctx := context.Background()

	// Arrange
	mockRepo := new(MockUserRepository)
	mockValidator := new(MockIdTokenValidator)
	mockEncoder := new(MockUserEncoder)
	passwordEncrypter := new(MockPasswordEncrypter)

	handler := user_application.NewGoogleSignInQueryHandler(mockRepo, mockValidator, mockEncoder, passwordEncrypter)

	// Act
	result, err := handler.Handle(ctx, nil)

	// Assert
	require.Nil(t, result)
	require.Error(t, err)
	require.Equal(t, "invalid query", err.Error())
}

func TestGoogleSignInQueryHandler_TokenValidationFailed(t *testing.T) {
	ctx := context.Background()

	// Arrange
	mockRepo := new(MockUserRepository)
	mockValidator := new(MockIdTokenValidator)
	mockEncoder := new(MockUserEncoder)
	passwordEncrypter := new(MockPasswordEncrypter)

	handler := user_application.NewGoogleSignInQueryHandler(mockRepo, mockValidator, mockEncoder, passwordEncrypter)

	idToken := "test-id-token"
	mockValidator.On("Validate", ctx, idToken).Return(nil, errors.New("invalid token"))

	// Act
	result, err := handler.Handle(ctx, &user_application.GoogleSignInQuery{IdToken: idToken})

	// Assert
	require.Nil(t, result)
	require.Error(t, err)
	require.Equal(t, "invalid token", err.Error())
	mockValidator.AssertCalled(t, "Validate", ctx, idToken)
}
