package user_ui

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	user_application "github.com/mik3lon/starter-template/internal/app/module/user/application"
	http_response "github.com/mik3lon/starter-template/internal/pkg/infrastructure/http/response"
	"github.com/mik3lon/starter-template/pkg/bus/command"
	"net/http"
)

type UpdateUserProfile struct {
	jw *http_response.JsonResponseWriter
	cb command.Bus
}

func NewUpdateUserProfile(
	cb command.Bus,
	jw *http_response.JsonResponseWriter,
) *UpdateUserProfile {
	return &UpdateUserProfile{cb: cb, jw: jw}
}

type UpdateProfileRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Surname  string `json:"surname" binding:"required"`
}

func (uup *UpdateUserProfile) HandleUpdateUserProfile(g *gin.Context) {
	email, exists := g.Get("user_email")
	if !exists {
		g.JSON(http.StatusBadRequest, gin.H{"error": errors.New("email not exists").Error()})
		return
	}

	var r UpdateProfileRequest

	if err := g.ShouldBindJSON(&r); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uup.cb.Dispatch(g, &user_application.UpdateUserProfileCommand{
		Email:    email.(string),
		Username: r.Username,
		Name:     r.Name,
		Surname:  r.Surname,
	})
	switch err.(type) {
	case nil:
		uup.jw.WriteResponse(g.Writer, "", http.StatusNoContent)
	default:
		fmt.Printf("error %v", err)
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	return
}
