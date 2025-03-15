package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"upload-gateway/internal/models"
)

type ServerRepository struct {
	db *gorm.DB
}

func NewServerRepository(db *gorm.DB) *ServerRepository {
	return &ServerRepository{
		db: db,
	}
}

func (r *ServerRepository) GetAvailableServers(ctx context.Context, limit int) (servers []*models.Server, err error) {
	query := r.db.Model(&models.Server{}).
		WithContext(ctx).
		Select("*").
		Limit(limit)

	if err = query.
		Find(&servers).Error; err != nil {
		err = fmt.Errorf("can't get servers: %w", err)
	}

	return
}
