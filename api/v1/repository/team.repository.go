package repository

import (
	"context"

	"github.com/backend-test-cubi-casa/api/v1/models"
	"gorm.io/gorm"
)

type ITeamRepository interface {
	Create(ctx context.Context, user *models.Team) error
	Find(ctx context.Context, params models.TeamSearchReq) ([]*models.Team, error)
}

type TeamRepository struct {
	db    *gorm.DB
	model *models.Team
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (repo TeamRepository) Create(ctx context.Context, model *models.Team) error {
	return repo.db.WithContext(ctx).Create(model).Error
}

func (repo TeamRepository) Find(ctx context.Context, params models.TeamSearchReq) ([]*models.Team, error) {
	var teams []*models.Team

	err := repo.
		buildingQuery(params).
		Preload(new(models.Hub).GetAlias()).
		WithContext(ctx).
		Find(&teams).
		Error

	return teams, err
}

func (repo TeamRepository) buildingQuery(params models.TeamSearchReq) *gorm.DB {
	queryBuilder := repo.db

	if params.ID > 0 {
		queryBuilder = queryBuilder.Where("id = ?", params.ID)
	}

	if params.Name != "" {
		queryBuilder = queryBuilder.Where("name LIKE ?", "%"+params.Name+"%")
	}

	return queryBuilder
}
