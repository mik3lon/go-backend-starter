package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	ur user_domain.UserRepository
	ue user_domain.UserEncoder
}

func NewAuthMiddleware(ur user_domain.UserRepository, ue user_domain.UserEncoder) *AuthMiddleware {
	return &AuthMiddleware{ur: ur, ue: ue}
}

// Check ensures that the user is authenticated
func (am *AuthMiddleware) Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		claims, err := am.ue.DecryptToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}

		mapClaims := claims.(jwt.MapClaims)
		c.Set("user_email", mapClaims["sub"])
	}
}
