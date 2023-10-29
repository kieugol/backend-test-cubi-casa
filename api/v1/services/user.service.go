package services

import (
	"context"
	"time"

	"github.com/backend-test-cubi-casa/api/v1/models"
	"github.com/backend-test-cubi-casa/api/v1/repository"
)

type IUserService interface {
	HandleCreate(req models.UserCreateReq) (*models.User, error)
	HandleSearch(req models.UserSearchReq) ([]*models.UserResp, error)
}

type UserService struct {
	repo repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (srv *UserService) HandleCreate(req models.UserCreateReq) (*models.User, error) {
	now := time.Now()
	user := &models.User{
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		TeamID:    req.TeamID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := srv.repo.Create(context.TODO(), user); err != nil {
		return nil, err
	}

	return user, nil
}

func (srv *UserService) HandleSearch(req models.UserSearchReq) ([]*models.UserResp, error) {
	users, err := srv.repo.Find(context.TODO(), req)
	if err != nil {
		return nil, err
	}

	var result []*models.UserResp
	for _, row := range users {
		result = append(result, srv.transformData(row))
	}

	return result, nil
}

func (srv *UserService) transformData(user *models.User) *models.UserResp {
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
