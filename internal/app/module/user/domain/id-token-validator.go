package user_domain

import "context"

type IdTokenClaims struct {
	Name              string
	Surname           string
	Username          string // Full name
	Email             string
	ProfilePictureUrl string
}

func NewIdTokenClaims(name string, surname string, username string, email string, profilePictureUrl string) *IdTokenClaims {
	return &IdTokenClaims{Name: name, Surname: surname, Username: username, Email: email, ProfilePictureUrl: profilePictureUrl}
}

type IdTokenValidator interface {
	Validate(ctx context.Context, idToken string) (*IdTokenClaims, error)
}
