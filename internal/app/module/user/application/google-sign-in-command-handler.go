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
	return "create-user-command"
}

type GoogleSignInQueryHandler struct {
	r  user_domain.UserRepository
	tv user_domain.IdTokenValidator
	ue user_domain.UserEncoder
}

func NewGoogleSignInQueryHandler(
	r user_domain.UserRepository,
	tv user_domain.IdTokenValidator,
	ue user_domain.UserEncoder,
) *GoogleSignInQueryHandler {
	return &GoogleSignInQueryHandler{r: r, tv: tv, ue: ue}
}

func (cuch GoogleSignInQueryHandler) Handle(ctx context.Context, c bus.Dto) (interface{}, error) {
	cuc, ok := c.(*GoogleSignInQuery)
	if !ok {
		return nil, errors.New("invalid command")
	}

	idTokenClaims, err := cuch.tv.Validate(ctx, cuc.IdToken)
	if err != nil {
		return nil, err
	}

	password, err := user_domain.GenerateHashedPassword(true, "")
	if err != nil {
		return nil, errors.New("failed to generate hashed password")
	}

	user := user_domain.CreateUser(
		uuid.NewString(),
		idTokenClaims.Username,
		idTokenClaims.Email,
		password,
		idTokenClaims.Name,
		idTokenClaims.Surname,
		"ROL_USER",
		idTokenClaims.ProfilePictureUrl,
	)

	err = cuch.r.Save(ctx, user)
	switch err.(type) {
	case nil, *user_domain.UserAlreadyExists:
		return cuch.ue.GenerateToken(user)
	default:
		return nil, err
	}
}
