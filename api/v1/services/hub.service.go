package services

import (
	"context"
	"time"

	"github.com/backend-test-cubi-casa/api/v1/models"
	"github.com/backend-test-cubi-casa/api/v1/repository"
)

type IHubService interface {
	HandleCreate(ctx context.Context, req models.HubCreateReq) (*models.Hub, error)
}

type HubService struct {
	repo repository.IHubRepository
}

func NewHubService(repo repository.IHubRepository) *HubService {
	return &HubService{
		repo: repo,
	}
}

func (srv *HubService) HandleCreate(ctx context.Context, req models.HubCreateReq) (*models.Hub, error) {
	now := time.Now()
	user := &models.Hub{
		Name:        req.Name,
		Location:    req.Location,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := srv.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
