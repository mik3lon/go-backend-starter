package user_application

import (
	"context"
	"errors"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"github.com/mik3lon/starter-template/pkg/bus"
)

type UserPasswordSignInQuery struct {
	Email    string
	Password string
}

func (c UserPasswordSignInQuery) Id() string {
	return "user-password-sign-in-query"
}

type UserPasswordSignInQueryHandler struct {
	r  user_domain.UserRepository
	ue user_domain.UserEncoder
	pe user_domain.PasswordEncrypter
}

func NewUserPasswordSignInQueryHandler(
	r user_domain.UserRepository,
	ue user_domain.UserEncoder,
	pe user_domain.PasswordEncrypter,
) *UserPasswordSignInQueryHandler {
	return &UserPasswordSignInQueryHandler{r: r, ue: ue, pe: pe}
}

func (upsq UserPasswordSignInQueryHandler) Handle(ctx context.Context, c bus.Dto) (interface{}, error) {
	cuc, ok := c.(*UserPasswordSignInQuery)
	if !ok {
		return nil, errors.New("invalid query")
	}

	user, err := upsq.r.FindByEmail(ctx, cuc.Email)
	if err != nil {
		return nil, err
	}

	err = upsq.pe.VerifyPassword(user.HashedPassword, cuc.Password)
	if err != nil {
		return nil, err
	}

	return upsq.ue.GenerateToken(user)
}
