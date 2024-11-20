package kernel

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib" // Import the pgx driver
	user_application "github.com/mik3lon/starter-template/internal/app/module/user/application"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	user_infrastructure "github.com/mik3lon/starter-template/internal/app/module/user/infrastructure"
	user_ui "github.com/mik3lon/starter-template/internal/app/module/user/ui"
	"github.com/mik3lon/starter-template/pkg/auth"
	"github.com/mik3lon/starter-template/pkg/config"
	"github.com/mik3lon/starter-template/pkg/http/middleware"
	"net/http"
)

const (
	GetUserList = "/users"
	GetUserMe   = "/users/me"
)

type UserModule struct {
	BaseModule

	UserSignInIndexHandler gin.HandlerFunc
	UserDashboardHandler   *user_ui.UserDashboardHandler

	GetUserList *user_ui.GetUserListHandler

	UserRepository            user_domain.UserRepository
	GoogleSocialSignInHandler *user_ui.GoogleSocialSignInHandler
	UserPasswordSignInHandler *user_ui.UserPasswordSignInHandler
	UserPasswordSignUpHandler *user_ui.UserPasswordSignUpHandler

	IdTokenValidator user_domain.IdTokenValidator
	GetUserMeHandler *user_ui.GetUserMeHandler
	UserEncoder      user_domain.UserEncoder
	AuthMiddleware   *middleware.AuthMiddleware
}

func (m *UserModule) Name() string {
	return "user_module"
}

// InitUserModule creates a new instance of NotificationModule.
func InitUserModule(k *Kernel, cnf *config.Config) *UserModule {
	_, err := user_infrastructure.NewPostgresUserRepository(cnf.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	r, err := user_infrastructure.NewPostgresUserRepository(cnf.DatabaseDSN)
	if err != nil {
		panic("error connecting with the database")
	}

	ue := auth.NewJWTUserEncoder(cnf.PrivateKeyPEM, cnf.PrivateKeyPassword, cnf.PublicKeyPEM)

	um := &UserModule{
		UserRepository:            r,
		UserEncoder:               ue,
		AuthMiddleware:            middleware.NewAuthMiddleware(r, ue),
		UserSignInIndexHandler:    user_ui.HandleUserSocialSignInIndex,
		UserDashboardHandler:      user_ui.NewUserDashboardHandler(k.JsonResponseWriter),
		GetUserList:               user_ui.NewGetUserListHandler(k.QueryBus, k.JsonResponseWriter),
		GoogleSocialSignInHandler: user_ui.NewGoogleSocialSignInHandler(k.QueryBus, k.JsonResponseWriter),
		IdTokenValidator:          user_infrastructure.NewGoogleIDTokenValidator(cnf.GoogleClientId),
		UserPasswordSignInHandler: user_ui.NewUserPasswordSignInHandler(k.QueryBus, k.JsonResponseWriter),
		UserPasswordSignUpHandler: user_ui.NewUserPasswordSignUpHandler(k.CommandBus, k.JsonResponseWriter),
		GetUserMeHandler:          user_ui.NewGetUserMeHandler(k.QueryBus, k.JsonResponseWriter),
	}

	pe := user_infrastructure.NewBcryptPasswordEncrypter()

	um.AddCommand(&user_application.CreateUserCommand{}, user_application.NewCreateUserCommandHandler(r, pe))

	um.AddQuery(&user_application.GoogleSignInQuery{}, user_application.NewGoogleSignInQueryHandler(r, um.IdTokenValidator, ue, pe))
	um.AddQuery(&user_application.FindUserQuery{}, user_application.NewFindUserQueryHandler(r))
	um.AddQuery(&user_application.ListUsersQuery{}, user_application.NewListUsersQueryHandler(r))
	um.AddQuery(&user_application.UserPasswordSignInQuery{}, user_application.NewUserPasswordSignInQueryHandler(r, ue, pe))

	return um
}

func (m *UserModule) RegisterRoutes(c *Kernel) {
	c.Router.Handle(
		http.MethodPost,
		"/users/social/signin/google",
		m.GoogleSocialSignInHandler.HandleGoogleSocialSignIn,
	)

	c.Router.Handle(
		http.MethodPost,
		"/users/auth/signin",
		m.UserPasswordSignInHandler.HandleUserPasswordSignIn,
	)

	c.Router.Handle(
		http.MethodPost,
		"/users/auth/signup",
		m.UserPasswordSignUpHandler.HandleUserPasswordSignUp,
	)

	c.Router.WithMiddleware().Handle(
		http.MethodGet,
		GetUserList,
		m.GetUserList.HandleGetUserList,
	)

	c.Router.WithMiddleware(m.AuthMiddleware.Check).Handle(
		http.MethodGet,
		GetUserMe,
		m.GetUserMeHandler.HandleGetUserMe,
	)
}
