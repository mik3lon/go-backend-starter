package user_ui

import (
	"errors"
	"github.com/gin-gonic/gin"
	user_application "github.com/mik3lon/starter-template/internal/app/module/user/application"
	http_response "github.com/mik3lon/starter-template/internal/pkg/infrastructure/http/response"
	"github.com/mik3lon/starter-template/pkg/bus/command"
	file2 "github.com/mik3lon/starter-template/pkg/file"
	"io"
	"net/http"
)

type UpdateUserProfilePhoto struct {
	jw *http_response.JsonResponseWriter
	cb command.Bus
}

func NewUpdateUserProfilePhoto(
	cb command.Bus,
	jw *http_response.JsonResponseWriter,
) *UpdateUserProfilePhoto {
	return &UpdateUserProfilePhoto{cb: cb, jw: jw}
}

func (uup *UpdateUserProfilePhoto) HandleUpdateProfilePhoto(g *gin.Context) {
	// Retrieve email from context
	email, exists := g.Get("user_email")
	if !exists {
		uup.jw.WriteErrorResponse(g.Writer, errors.New("email not exists"), http.StatusBadRequest, nil)
		return
	}

	file, err := g.FormFile("profile_image")
	if err != nil {
		uup.jw.WriteErrorResponse(g.Writer, err, http.StatusBadRequest, nil)
		return
	}

	f, err := file.Open()
	if err != nil {
		uup.jw.WriteErrorResponse(g.Writer, err, http.StatusInternalServerError, nil)
		return
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		uup.jw.WriteErrorResponse(g.Writer, err, http.StatusInternalServerError, nil)
		return
	}

	err = uup.cb.Dispatch(g, &user_application.UpdateUserProfilePhotoCommand{
		Email: email.(string),
		Image: file2.NewFileInfo(
			file.Filename,
			file.Header.Get("Content-Type"),
			file.Size,
			content,
		),
	})

	if err != nil {
		uup.jw.WriteErrorResponse(g.Writer, err, http.StatusInternalServerError, nil)
		return
	}

	uup.jw.WriteResponse(g.Writer, "", http.StatusNoContent)
}
