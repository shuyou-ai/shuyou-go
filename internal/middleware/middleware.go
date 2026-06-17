package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shuyou-ai/shuyou-go/internal/infra/jwt"
	"github.com/shuyou-ai/shuyou-go/pkg/response"
	"go.uber.org/zap"
)

func Auth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Unauthorized(c, "invalid authorization header")
			c.Abort()
			return
		}

		claims, err := jwtManager.Parse(parts[1])
		if err != nil {
			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		c.Set(string(UserIDKey), claims.UserID)
		c.Next()
	}
}

func Logger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", time.Since(start)),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
			log.Error("request failed", fields...)
			return
		}

		log.Info("request completed", fields...)
	}
}

func Recovery(log *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log.Error("panic recovered",
			zap.Any("error", recovered),
			zap.String("path", c.Request.URL.Path),
		)
		response.Fail(c, 500, 50000, "internal server error")
	})
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
