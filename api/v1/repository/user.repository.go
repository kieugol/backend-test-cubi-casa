package repository

import (
	"context"

	"github.com/backend-test-cubi-casa/api/v1/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(ctx context.Context, user *models.User) error
	Find(ctx context.Context, params models.UserSearchReq) ([]*models.User, error)
}

type UserRepository struct {
	db    *gorm.DB
	model *models.User
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo UserRepository) Create(ctx context.Context, model *models.User) error {
	return repo.db.WithContext(ctx).Create(model).Error
}

func (repo UserRepository) Find(ctx context.Context, params models.UserSearchReq) ([]*models.User, error) {
	var users []*models.User

	err := repo.
		buildingQuery(params).
		Preload("Team.Hub").
		WithContext(ctx).
		Find(&users).
		Error

	return users, err
}

func (repo UserRepository) buildingQuery(params models.UserSearchReq) *gorm.DB {
	queryBuilder := repo.db

	if params.Name != "" {
		queryBuilder = queryBuilder.Where("name LIKE ?", "%"+params.Name+"%")
	}

	if params.Email != "" {
		queryBuilder = queryBuilder.Where("email LIKE ?", "%"+params.Email+"%")
	}

	return queryBuilder
}
