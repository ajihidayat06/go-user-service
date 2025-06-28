package user

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct {
	// db *gorm.DB // contoh jika pakai GORM
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		// db: db,
	}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	// Implementasi simpan user ke DB
	return nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	// Implementasi cari user by email
	return nil, nil
}
