package models

import "time"

type UploadedFile struct {
	ID            uint64    `gorm:"id"`
	Ext           string    `gorm:"ext"`
	OriginalName  string    `gorm:"original_name"`
	GeneratedName string    `gorm:"generated_name"`
	CreatedAt     time.Time `gorm:"created_at"`
	UpdatedAt     time.Time `gorm:"updated_at"`

	UploadedFileChunks []*UploadedFileChunk
}

func (s *UploadedFile) TableName() string {
	return "uploaded_files"
}
