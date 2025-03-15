package internal

import (
	"context"
	"fmt"
	"upload-gateway/internal/models"
	"upload-gateway/internal/repository"
)

type DefaultServerSelector struct {
	serverRepository *repository.ServerRepository
}

const serverCount = 6

func NewDefaultServerSelector(serverRepository *repository.ServerRepository) *DefaultServerSelector {
	return &DefaultServerSelector{
		serverRepository: serverRepository,
	}
}

func (u *DefaultServerSelector) GetServers(ctx context.Context) ([]*models.Server, error) {
	servers, err := u.serverRepository.GetAvailableServers(ctx, serverCount)
	if err != nil {
		return nil, fmt.Errorf("error to get servers: %w", err)
	}

	return servers, nil
}
