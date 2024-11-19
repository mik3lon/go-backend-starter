package user_ui

import (
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	user_application "github.com/mik3lon/starter-template/internal/app/module/user/application"
	http_response "github.com/mik3lon/starter-template/internal/pkg/infrastructure/http/response"
	"github.com/mik3lon/starter-template/pkg/bus/query"
	"net/http"
	"strings"
)

type UserPasswordSignInHandler struct {
	jw *http_response.JsonResponseWriter
	qb query.Bus
}

func NewUserPasswordSignInHandler(
	qb query.Bus,
	jw *http_response.JsonResponseWriter,
) *UserPasswordSignInHandler {
	return &UserPasswordSignInHandler{qb: qb, jw: jw}
}

func (gss *UserPasswordSignInHandler) HandleUserPasswordSignIn(g *gin.Context) {
	email, password, ok := extractBasicAuth(g)
	if !ok {
		g.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("unauthorized").Error()})
		return
	}

	userToken, err := gss.qb.Ask(g, &user_application.UserPasswordSignInQuery{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return
	}

	switch err.(type) {
	case nil:
		gss.jw.WriteResponse(g.Writer, userToken, http.StatusOK)
	default:
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	return
}

func extractBasicAuth(c *gin.Context) (string, string, bool) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
		return "", "", false
	}

	// Decode the Base64 string
	encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		return "", "", false
	}

	// Split into username and password
	credentials := strings.SplitN(string(decodedBytes), ":", 2)
	if len(credentials) != 2 {
		return "", "", false
	}

	return credentials[0], credentials[1], true
}
