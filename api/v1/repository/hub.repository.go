package repository

import (
	"context"

	"github.com/backend-test-cubi-casa/api/v1/models"
	"gorm.io/gorm"
)

type IHubRepository interface {
	Create(ctx context.Context, user *models.Hub) error
}

type HubRepository struct {
	db    *gorm.DB
	model *models.Hub
}

func NewHubRepository(db *gorm.DB) *HubRepository {
	return &HubRepository{db: db}
}

func (repo HubRepository) Create(ctx context.Context, model *models.Hub) error {
	return repo.db.WithContext(ctx).Create(model).Error
}
