package transform

import "github.com/backend-test-cubi-casa/api/v1/models"

type TeamTransform struct {
}

func (srv *TeamTransform) TransformData(team *models.Team) *models.TeamResp {
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
