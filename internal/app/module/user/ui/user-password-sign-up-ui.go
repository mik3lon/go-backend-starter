package user_ui

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	user_application "github.com/mik3lon/starter-template/internal/app/module/user/application"
	http_response "github.com/mik3lon/starter-template/internal/pkg/infrastructure/http/response"
	"github.com/mik3lon/starter-template/pkg/bus/command"
	"net/http"
)

type UserPasswordSignUpHandler struct {
	jw *http_response.JsonResponseWriter
	cb command.Bus
}

type UserSignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func NewUserPasswordSignUpHandler(
	cb command.Bus,
	jw *http_response.JsonResponseWriter,
) *UserPasswordSignUpHandler {
	return &UserPasswordSignUpHandler{cb: cb, jw: jw}
}

func (gss *UserPasswordSignUpHandler) HandleUserPasswordSignUp(g *gin.Context) {
	var r UserSignUpRequest
	if err := g.ShouldBindJSON(&r); err != nil {
		gss.jw.WriteResponse(g.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	err := gss.cb.Dispatch(
		g,
		&user_application.CreateUserCommand{
			ID:                uuid.NewString(),
			Name:              r.Name,
			Surname:           "",
			Username:          r.Name,
			PlainPassword:     r.Password,
			Email:             r.Email,
			Role:              "ROLE_USER",
			ProfilePictureUrl: "",
			IsFormSocialAuth:  false,
		},
	)

	switch err.(type) {
	case nil:
		gss.jw.WriteResponse(g.Writer, "", http.StatusNoContent)
	default:
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	return
}
