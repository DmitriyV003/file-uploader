package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"upload-gateway/internal/models"
)

type UploadedFileChunkRepository struct {
	db *gorm.DB
}

func NewUploadedFileChunkRepository(db *gorm.DB) *UploadedFileChunkRepository {
	return &UploadedFileChunkRepository{
		db: db,
	}
}

func (r *UploadedFileChunkRepository) Save(ctx context.Context, uplFile *models.UploadedFileChunk) (err error) {
	if err = r.db.Save(uplFile).Error; err != nil {
		err = fmt.Errorf(
			"can't save uploaded file chunk: %w",
			err,
		)
	}

	return
}
