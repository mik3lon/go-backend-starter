package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"strings"
	"time"
)

type JWTUserEncoder struct {
	privateKeyPEM      string
	privateKeyPassword string
	publicKeyPEM       string
}

func NewJWTUserEncoder(privateKey string, password string, publicKey string) *JWTUserEncoder {
	return &JWTUserEncoder{
		privateKeyPEM:      strings.Replace(privateKey, `\n`, "\n", 100),
		privateKeyPassword: password,
		publicKeyPEM:       strings.Replace(publicKey, `\n`, "\n", 100),
	}
}

// loadPrivateKey loads an RSA private key from a PEM file with an optional privateKeyPassword
func loadPrivateKey(pemFile, password string) (*rsa.PrivateKey, error) {
	signBytes := []byte(pemFile)

	if password != "" {
		return jwt.ParseRSAPrivateKeyFromPEMWithPassword(signBytes, password)
	}

	return jwt.ParseRSAPrivateKeyFromPEM(signBytes)
}

// loadPublicKey loads an RSA public key from a PEM file
func loadPublicKey(pemFile string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemFile))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	// Parse the public key
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	// Assert the parsed key is an RSA public key
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return rsaPubKey, nil
}

// GenerateToken generates access and refresh tokens
func (jue *JWTUserEncoder) GenerateToken(user *user_domain.User) (*user_domain.TokenDetails, error) {
	// Load the private key
	privateKey, err := loadPrivateKey(jue.privateKeyPEM, jue.privateKeyPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %v", err)
	}

	// Set token expiration times
	accessTokenExpiration := time.Now().Add(2 * time.Hour).Unix()       // 15 minutes
	refreshTokenExpiration := time.Now().Add(7 * 24 * time.Hour).Unix() // 7 days

	// Create the access token
	accessClaims := jwt.MapClaims{
		"sub": user.Email,
		"exp": accessTokenExpiration,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	signedAccessToken, err := accessToken.SignedString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %v", err)
	}

	// Create the refresh token
	refreshClaims := jwt.MapClaims{
		"sub": user.Email,
		"exp": refreshTokenExpiration,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %v", err)
	}

	// Create and return TokenDetails
	tokenDetails := &user_domain.TokenDetails{
		UserEmail:           user.Email,
		AccessToken:         signedAccessToken,
		RefreshToken:        signedRefreshToken,
		AccessTokenExpires:  accessTokenExpiration,
		RefreshTokenExpires: refreshTokenExpiration,
	}
	return tokenDetails, nil
}

// DecryptToken verifies and parses a JWT token using the public key
func (jue *JWTUserEncoder) DecryptToken(tokenString string) (jwt.Claims, error) {
	// Load the public key
	publicKey, err := loadPublicKey(jue.publicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to load public key: %v", err)
	}

	// Parse and verify the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token.Claims, nil
}
