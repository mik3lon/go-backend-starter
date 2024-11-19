package user_domain

type PasswordEncrypter interface {
	GenerateHashedPassword(isSocial bool, plainPassword string) (string, error)
	VerifyPassword(hashedPassword, password string) error
}
