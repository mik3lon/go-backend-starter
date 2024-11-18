package user_infrastructure

import (
	"context"
	"fmt"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"google.golang.org/api/idtoken"
)

type GoogleIDTokenValidator struct {
	googleClientId string
}

func NewGoogleIDTokenValidator(googleClientId string) *GoogleIDTokenValidator {
	return &GoogleIDTokenValidator{googleClientId: googleClientId}
}

func (gv *GoogleIDTokenValidator) Validate(ctx context.Context, idToken string) (*user_domain.IdTokenClaims, error) {
	validator, err := idtoken.Validate(ctx, idToken, gv.googleClientId)
	if err != nil {
		return nil, fmt.Errorf("invalid ID token: %v", err)
	}

	payload := validator.Claims

	idTokenClaims := &user_domain.IdTokenClaims{
		Name:              payload["given_name"].(string),
		Email:             payload["email"].(string),
		ProfilePictureUrl: payload["picture"].(string),
	}

	if familyName, ok := payload["family_name"].(string); ok {
		idTokenClaims.Surname = familyName
	} else {
		idTokenClaims.Surname = ""
	}

	if fullName, ok := payload["name"].(string); ok {
		idTokenClaims.Username = fullName
	} else {
		idTokenClaims.Username = idTokenClaims.Name + " " + idTokenClaims.Surname
	}

	return idTokenClaims, nil
}
