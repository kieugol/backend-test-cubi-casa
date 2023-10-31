package services

import (
	"context"
	"time"

	"github.com/backend-test-cubi-casa/api/v1/models"
	"github.com/backend-test-cubi-casa/api/v1/repository"
	"github.com/backend-test-cubi-casa/api/v1/transforms"
)

type IUserService interface {
	HandleCreate(ctx context.Context, req models.UserCreateReq) (*models.User, error)
	HandleSearch(ctx context.Context, req models.UserSearchReq) ([]*models.UserResp, error)
}

type UserService struct {
	repo repository.IUserRepository
	trf  transform.UserTransform
}

func NewUserService(repo repository.IUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (srv *UserService) HandleCreate(ctx context.Context, req models.UserCreateReq) (*models.User, error) {
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

	if err := srv.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (srv *UserService) HandleSearch(ctx context.Context, req models.UserSearchReq) ([]*models.UserResp, error) {
	users, err := srv.repo.Find(ctx, req)
	if err != nil {
		return nil, err
	}

	var result []*models.UserResp
	for _, row := range users {
		result = append(result, srv.trf.TransformData(row))
	}

	return result, nil
}
