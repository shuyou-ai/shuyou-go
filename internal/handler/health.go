package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shuyou-ai/shuyou-go/pkg/response"
)

type HealthChecker interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	checker HealthChecker
}

func NewHealthHandler(checker HealthChecker) *HealthHandler {
	return &HealthHandler{checker: checker}
}

func (h *HealthHandler) Live(c *gin.Context) {
	response.Success(c, gin.H{"status": "up"})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	if err := h.checker.Ping(c.Request.Context()); err != nil {
		response.Fail(c, http.StatusServiceUnavailable, 50300, "database unavailable")
		return
	}

	response.Success(c, gin.H{"status": "ready"})
}
