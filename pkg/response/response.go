package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/shuyou-ai/shuyou-go/internal/errors"
)

type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{
		Code:    apperrors.CodeOK,
		Message: "ok",
		Data:    data,
	})
}

func Fail(c *gin.Context, httpStatus, code int, message string) {
	c.JSON(httpStatus, Body{
		Code:    code,
		Message: message,
	})
}

func Error(c *gin.Context, err error) {
	var appErr *apperrors.Error
	if errors.As(err, &appErr) {
		Fail(c, appErr.HTTPStatus, appErr.Code, appErr.Message)
		return
	}
	Fail(c, http.StatusInternalServerError, apperrors.CodeInternalError, apperrors.ErrInternalError.Message)
}

func BadRequest(c *gin.Context, message string) {
	Fail(c, http.StatusBadRequest, apperrors.CodeBadRequest, message)
}

func Unauthorized(c *gin.Context, message string) {
	Fail(c, http.StatusUnauthorized, apperrors.CodeUnauthorized, message)
}
