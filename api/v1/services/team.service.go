package services

import (
	"context"
	"time"

	"github.com/backend-test-cubi-casa/api/v1/models"
	"github.com/backend-test-cubi-casa/api/v1/repository"
)

type ITeamService interface {
	HandleCreate(req models.TeamCreateReq) (*models.Team, error)
	HandleSearch(req models.TeamSearchReq) ([]*models.TeamResp, error)
}

type TeamService struct {
	repo repository.ITeamRepository
}

func NewTeamService(repo repository.ITeamRepository) *TeamService {
	return &TeamService{
		repo: repo,
	}
}

func (srv *TeamService) HandleCreate(req models.TeamCreateReq) (*models.Team, error) {
	now := time.Now()
	user := &models.Team{
		Name:        req.Name,
		HubID:       req.HubID,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := srv.repo.Create(context.TODO(), user); err != nil {
		return nil, err
	}

	return user, nil
}

func (srv *TeamService) HandleSearch(req models.TeamSearchReq) ([]*models.TeamResp, error) {
	teams, err := srv.repo.Find(context.TODO(), req)
	if err != nil {
		return nil, err
	}

	var result []*models.TeamResp
	for _, row := range teams {
		result = append(result, srv.transformData(row))
	}

	return result, nil
}

func (srv *TeamService) transformData(team *models.Team) *models.TeamResp {
	teamResp := models.TeamResp{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		HubID:       team.HubID,
	}
	teamResp.Hub = models.HubResp{
		ID:          team.Hub.ID,
		Name:        team.Hub.Name,
		Location:    team.Hub.Location,
		Description: team.Hub.Description,
	}

	return &teamResp
}
