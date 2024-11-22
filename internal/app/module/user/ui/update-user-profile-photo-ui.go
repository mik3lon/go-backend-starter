package user_ui

import (
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
		g.JSON(http.StatusBadRequest, gin.H{"error": "email not exists"})
		return
	}

	// Retrieve the uploaded file
	file, err := g.FormFile("profile_image")
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Open the file
	f, err := file.Open()
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	// Read file content
	content, err := io.ReadAll(f)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Dispatch the command
	err = uup.cb.Dispatch(g, &user_application.UpdateUserProfilePhotoCommand{
		Email: email.(string),
		Image: file2.NewFileInfo(
			file.Filename,
			file.Header.Get("Content-Type"),
			file.Size,
			content,
		),
	})

	// Handle response based on the result of dispatch
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Use 204 No Content for successful execution
	g.Status(http.StatusNoContent)
}
