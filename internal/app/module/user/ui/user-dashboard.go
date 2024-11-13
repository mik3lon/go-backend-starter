package user_ui

import (
	"github.com/gin-gonic/gin"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	http_response "github.com/mik3lon/starter-template/internal/pkg/infrastructure/http/response"
	"log"
	"net/http"
)

type UserDashboardHandler struct {
	jrw *http_response.JsonResponseWriter
}

func NewUserDashboardHandler(jrw *http_response.JsonResponseWriter) *UserDashboardHandler {
	return &UserDashboardHandler{jrw: jrw}
}

func (uth *UserDashboardHandler) HandleUserDashboard(g *gin.Context) {
	u, exists := g.Get("user")
	if !exists {
		log.Println("UserInfo not found in context")
		g.JSON(http.StatusInternalServerError, gin.H{"error": "User info not found"})
		return
	}

	user, ok := u.(user_domain.User)
	if !ok {
		panic("user not found")
	}

	uth.jrw.WriteResponse(g.Writer, user, http.StatusOK)
}
