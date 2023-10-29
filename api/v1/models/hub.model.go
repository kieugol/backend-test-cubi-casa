package models

import (
	"time"

	"github.com/backend-test-cubi-casa/helpers/util"
)

func (m *Hub) GetAlias() string {
	return util.GetStructName(Hub{})
}

type Hub struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Location    string    `gorm:"type:varchar(255);not null" json:"location"`
	Description string    `gorm:"type:varchar(255);" json:"description"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at"`
}

type HubCreateReq struct {
	Name        string `json:"name" binding:"required"`
	Location    string `json:"location" binding:"required"`
	Description string `json:"description" binding:"max=255"`
}

type HubResp struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
}
