package models

type Server struct {
	ID  uint64 `gorm:"id"`
	URL string `gorm:"url"`
}

func (s *Server) TableName() string {
	return "servers"
}
