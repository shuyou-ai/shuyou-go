package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shuyou-ai/shuyou-go/internal/dto"
	apperrors "github.com/shuyou-ai/shuyou-go/internal/errors"
	"github.com/shuyou-ai/shuyou-go/internal/infra/jwt"
	"github.com/shuyou-ai/shuyou-go/internal/middleware"
	"github.com/shuyou-ai/shuyou-go/internal/service"
	"github.com/shuyou-ai/shuyou-go/pkg/response"
	pkgvalidator "github.com/shuyou-ai/shuyou-go/pkg/validator"
)

type UserHandler struct {
	userService service.UserService
	jwtManager  *jwt.Manager
}

func NewUserHandler(userService service.UserService, jwtManager *jwt.Manager) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtManager:  jwtManager,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, pkgvalidator.FormatErrors(err))
		return
	}

	user, err := h.userService.Register(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		response.Error(c, err)
		return
	}

	token, err := h.jwtManager.Generate(user.ID)
	if err != nil {
		response.Error(c, apperrors.Wrap(err, apperrors.CodeInternalError, "generate token failed"))
		return
	}

	response.Success(c, dto.AuthResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, pkgvalidator.FormatErrors(err))
		return
	}

	user, err := h.userService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Error(c, err)
		return
	}

	token, err := h.jwtManager.Generate(user.ID)
	if err != nil {
		response.Error(c, apperrors.Wrap(err, apperrors.CodeInternalError, "generate token failed"))
		return
	}

	response.Success(c, dto.AuthResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString(string(middleware.UserIDKey))
	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.ToUserResponse(user))
}
