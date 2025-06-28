// internal/pkg/response/response.go
package response

import (
	"net/http"

	"go-user-service/internal/pkg/errors"

	"github.com/gin-gonic/gin"
)

// Response represents standardized API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// Success responses
func JSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

func JSONWithMeta(c *gin.Context, statusCode int, data interface{}, meta *Meta) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Error responses
func Error(c *gin.Context, err error) {
	var appErr *errors.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.StatusCode, Response{
			Success: false,
			Error: &ErrorInfo{
				Code:    string(appErr.Code),
				Message: appErr.Message,
				Details: appErr.Details,
			},
		})
		return
	}

	// Handle unknown errors
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    string(errors.ErrCodeInternal),
			Message: "Internal server error",
		},
	})
}

// Convenience methods
func OK(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, data)
}

func Created(c *gin.Context, data interface{}) {
	JSON(c, http.StatusCreated, data)
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func BadRequest(c *gin.Context, message string) {
	err := errors.New(errors.ErrCodeValidation, message)
	Error(c, err)
}

func NotFound(c *gin.Context, message string) {
	err := errors.New(errors.ErrCodeNotFound, message)
	Error(c, err)
}

func Unauthorized(c *gin.Context, message string) {
	err := errors.New(errors.ErrCodeUnauthorized, message)
	Error(c, err)
}

func InternalError(c *gin.Context, err error) {
	appErr := errors.Wrap(err, errors.ErrCodeInternal, "Internal server error")
	Error(c, appErr)
}
