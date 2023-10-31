package services

import (
	"context"
	"time"

	"github.com/backend-test-cubi-casa/api/v1/models"
	"github.com/backend-test-cubi-casa/api/v1/repository"
	transform "github.com/backend-test-cubi-casa/api/v1/transforms"
)

type ITeamService interface {
	HandleCreate(ctx context.Context, req models.TeamCreateReq) (*models.Team, error)
	HandleSearch(ctx context.Context, req models.TeamSearchReq) ([]*models.TeamResp, error)
}

type TeamService struct {
	repo repository.ITeamRepository
	trf  transform.TeamTransform
}

func NewTeamService(repo repository.ITeamRepository) *TeamService {
	return &TeamService{
		repo: repo,
	}
}

func (srv *TeamService) HandleCreate(ctx context.Context, req models.TeamCreateReq) (*models.Team, error) {
	now := time.Now()
	user := &models.Team{
		Name:        req.Name,
		HubID:       req.HubID,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := srv.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (srv *TeamService) HandleSearch(ctx context.Context, req models.TeamSearchReq) ([]*models.TeamResp, error) {
	teams, err := srv.repo.Find(ctx, req)
	if err != nil {
		return nil, err
	}

	var result []*models.TeamResp
	for _, row := range teams {
		result = append(result, srv.trf.TransformData(row))
	}

	return result, nil
}
