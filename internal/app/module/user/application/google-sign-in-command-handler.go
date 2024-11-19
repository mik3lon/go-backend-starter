package user_application

import (
	"context"
	"errors"
	"github.com/google/uuid"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"github.com/mik3lon/starter-template/pkg/bus"
)

type GoogleSignInQuery struct {
	IdToken string
}

func (c GoogleSignInQuery) Id() string {
	return "google-sign-in-query"
}

type GoogleSignInQueryHandler struct {
	r  user_domain.UserRepository
	tv user_domain.IdTokenValidator
	ue user_domain.UserEncoder
	pe user_domain.PasswordEncrypter
}

func NewGoogleSignInQueryHandler(
	r user_domain.UserRepository,
	tv user_domain.IdTokenValidator,
	ue user_domain.UserEncoder,
	pe user_domain.PasswordEncrypter,
) *GoogleSignInQueryHandler {
	return &GoogleSignInQueryHandler{r: r, tv: tv, ue: ue, pe: pe}
}

func (cuch GoogleSignInQueryHandler) Handle(ctx context.Context, c bus.Dto) (interface{}, error) {
	cuc, ok := c.(*GoogleSignInQuery)
	if !ok {
		return nil, errors.New("invalid query")
	}

	idTokenClaims, err := cuch.tv.Validate(ctx, cuc.IdToken)
	if err != nil {
		return nil, err
	}

	user, err := cuch.r.FindByEmail(ctx, idTokenClaims.Email)
	switch {
	case err == nil:
	case errors.As(err, new(*user_domain.UserNotFound)):
		// User not found, create a new one
		password, genErr := cuch.pe.GenerateHashedPassword(true, "")
		if genErr != nil {
			return nil, errors.New("failed to generate hashed password")
		}

		user = user_domain.CreateUser(
			uuid.NewString(),
			idTokenClaims.Username,
			idTokenClaims.Email,
			password,
			idTokenClaims.Name,
			idTokenClaims.Surname,
			"ROL_USER",
			idTokenClaims.ProfilePictureUrl,
		)

		saveErr := cuch.r.Save(ctx, user)
		if saveErr != nil {
			return nil, saveErr
		}
	default:
		return nil, err
	}

	return cuch.ue.GenerateToken(user)
}
