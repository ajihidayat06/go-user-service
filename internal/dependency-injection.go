package internal

import (
	"go-user-service/internal/user"

	"gorm.io/gorm"
)

func diUser(db *gorm.DB) *user.Handler {
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	return userHandler
}
