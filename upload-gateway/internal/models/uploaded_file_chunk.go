package models

import "time"

type UploadedFileChunk struct {
	ID             uint64    `gorm:"id"`
	UploadedFileID uint64    `gorm:"uploaded_file_id"`
	ChunkNumber    uint      `gorm:"chunk_number"`
	Name           string    `gorm:"name"`
	Size           uint64    `gorm:"size"`
	ServerID       uint64    `gorm:"server_id"`
	Hash           string    `gorm:"hash"`
	CreatedAt      time.Time `gorm:"created_at"`
	UpdatedAt      time.Time `gorm:"updated_at"`

	Server *Server
}

func (s *UploadedFileChunk) TableName() string {
	return "uploaded_file_chunks"
}
