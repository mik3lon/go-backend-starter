package kernel

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	http_response "github.com/mik3lon/starter-template/internal/pkg/infrastructure/http/response"
	"github.com/mik3lon/starter-template/pkg/bus/command"
	"github.com/mik3lon/starter-template/pkg/bus/query"
	"github.com/mik3lon/starter-template/pkg/config"
	"github.com/mik3lon/starter-template/pkg/file"
	"github.com/mik3lon/starter-template/pkg/http/middleware"
	shared_image_infrastructure "github.com/mik3lon/starter-template/pkg/infrastructure"
	"github.com/mik3lon/starter-template/pkg/router"
	"net/http"
)

type Kernel struct {
	Router             *router.GinRouter
	Modules            map[string]Module
	server             *http.Server
	CommandBus         *command.CommandBus
	QueryBus           *query.QueryBus
	JsonResponseWriter *http_response.JsonResponseWriter

	AuthMiddleware *middleware.AuthMiddleware
	ImageUploader  file.ImageUploader
}

// Init initializes the container with a router implementation.
func Init(cnf *config.Config) *Kernel {
	r := router.NewGinRouter()

	l := shared_image_infrastructure.NewZerologAdapter()

	k := &Kernel{
		Router: r,
		server: &http.Server{
			Addr:    cnf.AddressPort,
			Handler: r.Handler(),
		},
		CommandBus:         command.InitCommandBus(l),
		QueryBus:           query.InitQueryBus(l),
		JsonResponseWriter: http_response.NewJsonResponseWriter(),
		ImageUploader:      buildImageUploader(buildS3Client(cnf), cnf, l),
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

func buildS3Endpoint(cnf *config.Config) string {
	if cnf.AppEnv == "test" {
		return "http://localhost:4566"
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com", cnf.S3ImageBucket, cnf.S3Region)
}

func buildImageUploader(client *s3.S3, cnf *config.Config, l shared_image_infrastructure.Logger) file.ImageUploader {
	s3Endpoint := buildS3Endpoint(cnf)
	return shared_image_infrastructure.NewS3ImageUploader(client, cnf.S3ImageBucket, s3Endpoint, l)
}

func buildS3Client(config *config.Config) *s3.S3 {
	if config.AppEnv == "test" {
		sess, err := session.NewSession(&aws.Config{
			Region:           aws.String(config.S3Region),
			Endpoint:         aws.String(config.S3Endpoint),
			S3ForcePathStyle: aws.Bool(true),
		})

		if err != nil {
			panic(err)
		}

		return s3.New(sess)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String(config.S3Region)},
	}))

	return s3.New(sess)
}
