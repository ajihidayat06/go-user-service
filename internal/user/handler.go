package user

import (
	"go-user-service/internal/pkg/errors"
	"go-user-service/internal/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.GET("/", h.GetAll)
	users.POST("/register", h.Register)
}

func (h *Handler) GetAll(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "list all users"})
}

func (h *Handler) Register(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := errors.New(errors.ErrCodeValidation, "Invalid user ID format")
		response.Error(c, err)
		return
	}
	user, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, user)
}
