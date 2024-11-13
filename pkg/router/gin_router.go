package router

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"net/http"
)

type Middleware func() gin.HandlerFunc

type GinRouter struct {
	engine     *gin.Engine
	middleware []Middleware
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	Name          string `json:"name"`        // Full name
	GivenName     string `json:"given_name"`  // First name
	FamilyName    string `json:"family_name"` // Last name or surname
}

func NewGinRouter() *GinRouter {
	gob.Register(user_domain.User{})
	engine := gin.New()

	return &GinRouter{
		engine: engine,
	}
}

func (g *GinRouter) Handle(method, path string, handler gin.HandlerFunc) {
	// Apply middleware to the handler
	finalHandler := handler
	for _, m := range g.middleware {
		finalHandler = wrapMiddleware(m, finalHandler)
	}

	// Register the route
	switch method {
	case http.MethodGet:
		g.engine.GET(path, finalHandler)
	case http.MethodPost:
		g.engine.POST(path, finalHandler)
	case http.MethodPut:
		g.engine.PUT(path, finalHandler)
	case http.MethodDelete:
		g.engine.DELETE(path, finalHandler)
	default:
		// Handle other HTTP methods if needed
	}
}

func (g *GinRouter) WithMiddleware(middleware ...Middleware) *GinRouter {
	g.middleware = append(g.middleware, middleware...)
	return g
}

func (g *GinRouter) Serve(addr string) error {
	return g.engine.Run(addr)
}

func (g *GinRouter) Handler() http.Handler {
	return g.engine
}

// Helper function to wrap Middleware to gin.HandlerFunc
func wrapMiddleware(m Middleware, next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		m()(c) // Execute the middleware
		if c.IsAborted() {
			return // If middleware aborted the request, stop the chain
		}
		next(c) // Proceed to the next handler
	}
}
