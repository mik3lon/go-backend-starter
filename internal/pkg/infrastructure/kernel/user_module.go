package kernel

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib" // Import the pgx driver
	user_application "github.com/mik3lon/starter-template/internal/app/module/user/application"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	user_infrastructure "github.com/mik3lon/starter-template/internal/app/module/user/infrastructure"
	user_ui "github.com/mik3lon/starter-template/internal/app/module/user/ui"
	"github.com/mik3lon/starter-template/pkg/config"
	"net/http"
)

const (
	GetUserList = "/users"
)

type UserModule struct {
	BaseModule

	UserSignInIndexHandler gin.HandlerFunc
	UserDashboardHandler   *user_ui.UserDashboardHandler

	GetUserList *user_ui.GetUserListHandler

	UserRepository user_domain.UserRepository
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

	um := &UserModule{
		UserRepository:         r,
		UserSignInIndexHandler: user_ui.HandleUserSocialSignInIndex,
		UserDashboardHandler:   user_ui.NewUserDashboardHandler(k.JsonResponseWriter),
		GetUserList:            user_ui.NewGetUserListHandler(k.QueryBus, k.JsonResponseWriter),
	}

	um.AddCommand(&user_application.CreateUserCommand{}, user_application.NewCreateUserCommandHandler(r))
	um.AddQuery(&user_application.FindUserQuery{}, user_application.NewFindUserQueryHandler(r))
	um.AddQuery(&user_application.ListUsersQuery{}, user_application.NewListUsersQueryHandler(r))

	return um
}

func (m *UserModule) RegisterRoutes(c *Kernel) {

	c.Router.WithMiddleware().Handle(
		http.MethodGet,
		"/dashboard",
		m.UserDashboardHandler.HandleUserDashboard,
	)

	c.Router.WithMiddleware().Handle(
		http.MethodGet,
		GetUserList,
		m.GetUserList.HandleGetUserList,
	)
}
