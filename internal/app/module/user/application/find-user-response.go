package user_application

import (
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"time"
)

type FindUserResponse struct {
	ID                string    `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Name              string    `json:"name"`
	Surname           string    `json:"surname"`
	Role              string    `json:"role"`
	ProfilePictureUrl string    `json:"profile_picture_url"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func NewFindUserResponseFromUser(u *user_domain.User) *FindUserResponse {
	return &FindUserResponse{
		ID:                u.ID,
		Username:          u.Username,
		Email:             u.Email,
		Name:              u.Name,
		Surname:           u.Surname,
		Role:              u.Role,
		ProfilePictureUrl: u.ProfilePictureUrl,
		CreatedAt:         u.CreatedAt,
		UpdatedAt:         u.UpdatedAt,
	}
}
