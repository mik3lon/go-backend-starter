package kernel

import (
	"context"
	http_response "github.com/mik3lon/starter-template/internal/pkg/infrastructure/http/response"
	"github.com/mik3lon/starter-template/pkg/bus/command"
	"github.com/mik3lon/starter-template/pkg/bus/query"
	"github.com/mik3lon/starter-template/pkg/config"
	"github.com/mik3lon/starter-template/pkg/http/middleware"
	"github.com/mik3lon/starter-template/pkg/router"
	"github.com/rs/zerolog"
	"net/http"
	"os"
)

type Kernel struct {
	Router             *router.GinRouter
	Modules            map[string]Module
	server             *http.Server
	CommandBus         *command.CommandBus
	QueryBus           *query.QueryBus
	JsonResponseWriter *http_response.JsonResponseWriter

	AuthMiddleware *middleware.AuthMiddleware
}

// Init initializes the container with a router implementation.
func Init(cnf *config.Config) *Kernel {
	r := router.NewGinRouter()

	l := zerolog.New(os.Stderr).With().Timestamp().Logger()
	k := &Kernel{
		Router: r,
		server: &http.Server{
			Addr:    cnf.AddressPort,
			Handler: r.Handler(),
		},
		CommandBus:         command.InitCommandBus(l),
		QueryBus:           query.InitQueryBus(l),
		JsonResponseWriter: http_response.NewJsonResponseWriter(),
	}

	userModule := InitUserModule(k, cnf)
	k.addModule(userModule)

	k.RegisterModuleRoutes()

	return k
}

// RegisterModuleRoutes allows each module to register its routes.
func (k *Kernel) RegisterModuleRoutes() {
	for _, m := range k.Modules {
		m.RegisterRoutes(k)
	}
}

// StartServer starts the HTTP server.
func (k *Kernel) StartServer() error {
	return k.server.ListenAndServe()
}

func (k *Kernel) addModule(module Module) {
	if k.Modules == nil {
		k.Modules = make(map[string]Module)
	}

	if k.Modules[module.Name()] != nil {
		panic("Module already exists")
	}
	k.Modules[module.Name()] = module

	for c, ch := range module.Commands() {
		err := k.CommandBus.RegisterCommand(c, ch)
		if err != nil {
			panic(err)
		}
	}

	for q, ch := range module.Queries() {
		err := k.QueryBus.RegisterQuery(q, ch)
		if err != nil {
			panic(err)
		}
	}
}

func (k *Kernel) ShutdownServer(ctx context.Context) error {
	return k.server.Shutdown(ctx)
}
