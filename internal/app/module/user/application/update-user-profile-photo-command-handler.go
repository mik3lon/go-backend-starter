package user_application

import (
	"context"
	"errors"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"github.com/mik3lon/starter-template/pkg/bus"
	"github.com/mik3lon/starter-template/pkg/file"
)

type UpdateUserProfilePhotoCommand struct {
	Email string
	Image *file.FileInfo
}

func (c UpdateUserProfilePhotoCommand) Id() string {
	return "find-user-query-handler"
}

type UpdateUserProfilePhotoCommandHandler struct {
	r  user_domain.UserRepository
	iu file.ImageUploader
}

func NewUpdateUserProfilePhotoCommandHandler(r user_domain.UserRepository, iu file.ImageUploader) *UpdateUserProfilePhotoCommandHandler {
	return &UpdateUserProfilePhotoCommandHandler{r: r, iu: iu}
}

func (uupch UpdateUserProfilePhotoCommandHandler) Handle(ctx context.Context, command bus.Dto) error {
	c, ok := command.(*UpdateUserProfilePhotoCommand)
	if !ok {
		return errors.New("invalid command")
	}

	user, err := uupch.r.FindByEmail(ctx, c.Email)
	if err != nil {
		return err
	}

	upload, err := uupch.iu.Upload(ctx, *c.Image)
	if err != nil {
		return err
	}

	user.UpdateProfilePhoto(upload.Url)

	return uupch.r.Save(ctx, user)
}
