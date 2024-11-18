package user_domain

type UserEncoder interface {
	GenerateToken(user *User) (*TokenDetails, error)
}
