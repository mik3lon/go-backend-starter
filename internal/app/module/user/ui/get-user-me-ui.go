package user_ui

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	user_application "github.com/mik3lon/starter-template/internal/app/module/user/application"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	http_response "github.com/mik3lon/starter-template/internal/pkg/infrastructure/http/response"
	"github.com/mik3lon/starter-template/pkg/bus/query"
	"net/http"
)

type GetUserMeHandler struct {
	jw *http_response.JsonResponseWriter
	qb query.Bus
}

func NewGetUserMeHandler(
	qb query.Bus,
	jw *http_response.JsonResponseWriter,
) *GetUserMeHandler {
	return &GetUserMeHandler{qb: qb, jw: jw}
}

func (gss *GetUserMeHandler) HandleGetUserMe(g *gin.Context) {
	email, exists := g.Get("user_email")
	if !exists {
		g.JSON(http.StatusBadRequest, gin.H{"error": errors.New("email not exists").Error()})
		return
	}

	userResponse, err := gss.qb.Ask(g, &user_application.FindUserQuery{Email: email.(string)})
	switch err.(type) {
	case nil:
		gss.jw.WriteResponse(g.Writer, userResponse, http.StatusOK)
	case *user_domain.UserNotFound:
		errs := make([]error, 0)
		errs = append(errs, err)
		gss.jw.WriteErrorResponse(g.Writer, errs, http.StatusNotFound, err)
	default:
		fmt.Printf("error %v", err)
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	return
}
