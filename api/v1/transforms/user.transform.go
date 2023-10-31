package transform

import "github.com/backend-test-cubi-casa/api/v1/models"

type UserTransform struct {
}

func (trf UserTransform) TransformData(user *models.User) *models.UserResp {
	userResp := models.UserResp{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Phone:   user.Phone,
		Address: user.Address,
		TeamID:  user.TeamID,
	}
	userResp.Team = models.TeamResp{
		ID:          user.Team.ID,
		Name:        user.Team.Name,
		HubID:       user.Team.HubID,
		Description: user.Team.Description,
	}
	userResp.Team.Hub = models.HubResp{
		ID:          user.Team.Hub.ID,
		Name:        user.Team.Hub.Name,
		Location:    user.Team.Hub.Location,
		Description: user.Team.Hub.Description,
	}

	return &userResp
}
