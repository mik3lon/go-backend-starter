package user_ui

import (
	"github.com/gin-gonic/gin"
	user_application "github.com/mik3lon/starter-template/internal/app/module/user/application"
	http_response "github.com/mik3lon/starter-template/internal/pkg/infrastructure/http/response"
	"github.com/mik3lon/starter-template/pkg/bus/query"
	"net/http"
	"strconv"
)

type GetUserListHandler struct {
	qb  query.Bus
	jrw *http_response.JsonResponseWriter
}

func NewGetUserListHandler(qb query.Bus, jrw *http_response.JsonResponseWriter) *GetUserListHandler {
	return &GetUserListHandler{qb: qb, jrw: jrw}
}

func (uth *GetUserListHandler) HandleGetUserList(g *gin.Context) {
	page, err := strconv.Atoi(g.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(g.DefaultQuery("size", "20"))
	if err != nil || size < 1 {
		size = 20
	}

	userList, err := uth.qb.Ask(g, &user_application.ListUsersQuery{
		Page:   page,
		Size:   size,
		Filter: nil,
	})
	if err != nil {
		panic(err)
	}

	uth.jrw.WriteResponse(g.Writer, userList, http.StatusOK)
}
