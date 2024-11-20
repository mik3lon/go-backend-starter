package user_ui

import (
	"fmt"
	"github.com/gin-gonic/gin"
	user_application "github.com/mik3lon/starter-template/internal/app/module/user/application"
	http_response "github.com/mik3lon/starter-template/internal/pkg/infrastructure/http/response"
	"github.com/mik3lon/starter-template/pkg/bus/query"
	"net/http"
)

type GoogleSigInRequest struct {
	IdToken string `json:"id_token"`
}

type GoogleSocialSignInHandler struct {
	jw *http_response.JsonResponseWriter
	qb query.Bus
}

func NewGoogleSocialSignInHandler(
	qb query.Bus,
	jw *http_response.JsonResponseWriter,
) *GoogleSocialSignInHandler {
	return &GoogleSocialSignInHandler{qb: qb, jw: jw}
}

func (gss *GoogleSocialSignInHandler) HandleGoogleSocialSignIn(g *gin.Context) {
	var r GoogleSigInRequest

	if err := g.ShouldBindJSON(&r); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userToken, err := gss.qb.Ask(g, &user_application.GoogleSignInQuery{IdToken: r.IdToken})
	switch err.(type) {
	case nil:
		gss.jw.WriteResponse(g.Writer, userToken, http.StatusOK)
	default:
		fmt.Printf("error %v", err)
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	return
}
