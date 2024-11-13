package user_ui

import (
	"github.com/gin-gonic/gin"
)

func HandleUserSocialSignInIndex(g *gin.Context) {
	g.Writer.Header().Set("Content-Type", "text/html")
}
