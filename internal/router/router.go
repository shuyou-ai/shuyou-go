package router

import (
	"github.com/gin-gonic/gin"
	"github.com/shuyou-ai/shuyou-go/internal/handler"
	"github.com/shuyou-ai/shuyou-go/internal/infra/jwt"
	"github.com/shuyou-ai/shuyou-go/internal/middleware"
	"go.uber.org/zap"
)

type Handlers struct {
	Health *handler.HealthHandler
	User   *handler.UserHandler
}

type Deps struct {
	Log      *zap.Logger
	JWT      *jwt.Manager
	Handlers *Handlers
}

func New(mode string, deps Deps) *gin.Engine {
	gin.SetMode(mode)

	r := gin.New()
	r.Use(
		middleware.Recovery(deps.Log),
		middleware.Logger(deps.Log),
		middleware.CORS(),
	)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health/live", deps.Handlers.Health.Live)
		v1.GET("/health/ready", deps.Handlers.Health.Ready)

		users := v1.Group("/users")
		{
			users.POST("/register", deps.Handlers.User.Register)
			users.POST("/login", deps.Handlers.User.Login)

			auth := users.Group("")
			auth.Use(middleware.Auth(deps.JWT))
			auth.GET("/me", deps.Handlers.User.GetProfile)
		}
	}

	return r
}
