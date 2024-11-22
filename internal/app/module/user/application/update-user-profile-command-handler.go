package user_application

import (
	"context"
	"errors"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"github.com/mik3lon/starter-template/pkg/bus"
)

type UpdateUserProfileCommand struct {
	Email    string
	Username string
	Name     string
	Surname  string
}

func (c UpdateUserProfileCommand) Id() string {
	return "find-user-query-handler"
}

type UpdateUserProfileCommandHandler struct {
	r user_domain.UserRepository
}

func NewUpdateUserProfileCommandHandler(r user_domain.UserRepository) *UpdateUserProfileCommandHandler {
	return &UpdateUserProfileCommandHandler{r: r}
}

func (uupch UpdateUserProfileCommandHandler) Handle(ctx context.Context, command bus.Dto) error {
	c, ok := command.(*UpdateUserProfileCommand)
	if !ok {
		return errors.New("invalid command")
	}

	user, err := uupch.r.FindByEmail(ctx, c.Email)
	if err != nil {
		return err
	}

	user.UpdateProfile(c.Username, c.Name, c.Surname)

	return uupch.r.Save(ctx, user)
}
