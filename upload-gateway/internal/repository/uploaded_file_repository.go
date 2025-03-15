package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"upload-gateway/internal/models"
)

type UploadedFileRepository struct {
	db *gorm.DB
}

func NewUploadedFileRepository(db *gorm.DB) *UploadedFileRepository {
	return &UploadedFileRepository{
		db: db,
	}
}

func (r *UploadedFileRepository) Save(ctx context.Context, uplFile *models.UploadedFile) (err error) {
	if err = r.db.Save(uplFile).Error; err != nil {
		err = fmt.Errorf(
			"can't save uploaded file: %w",
			err,
		)
	}

	return
}

func (r *UploadedFileRepository) GetByGeneratedName(ctx context.Context, name string) (uf *models.UploadedFile, err error) {
	if err = r.db.
		WithContext(ctx).
		Where("generated_name = ?", name).
		Preload("UploadedFileChunks", func(db *gorm.DB) *gorm.DB {
			return db.Order("uploaded_file_chunks.chunk_number")
		}).
		Preload("UploadedFileChunks.Server").
		First(&uf).Error; err != nil {
		return nil, err
	}

	return
}
