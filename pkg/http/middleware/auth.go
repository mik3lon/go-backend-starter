package middleware

import (
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware ensures that the user is authenticated
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")

		if user == nil {
			// If user is not authenticated, respond with an error
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		_, ok := user.(user_domain.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid session data"})
			c.Abort()
			return
		}

		// Store GoogleUserInfo in the context
		c.Set("user", user)

		c.Next()
	}
}
