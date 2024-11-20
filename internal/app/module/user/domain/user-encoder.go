package user_domain

import "github.com/golang-jwt/jwt"

type UserEncoder interface {
	GenerateToken(user *User) (*TokenDetails, error)
	DecryptToken(tokenString string) (jwt.Claims, error)
}
