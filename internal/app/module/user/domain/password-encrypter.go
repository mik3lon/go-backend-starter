package user_domain

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", errors.New("error generating random token")
	}
	// Encode the token to base64 for safe storage
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func GenerateHashedPassword(isSocial bool, plainPassword string) (string, error) {
	if isSocial {
		return generateRandomToken(32)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("error generating password hash")
	}
	return string(hash), nil
}
