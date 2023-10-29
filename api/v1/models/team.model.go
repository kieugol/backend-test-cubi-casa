package models

import (
	"time"

	"github.com/backend-test-cubi-casa/helpers/util"
)

func (m *Team) GetAlias() string {
	return util.GetStructName(Team{})
}

type Team struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:varchar(255);" json:"description"`
	HubID       uint      `json:"hub_id"`
	Hub         Hub       `gorm:"foreignKey:HubID" json:"-"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at"`
}

type TeamCreateReq struct {
	Name        string `json:"name" binding:"required,max=255"`
	HubID       uint   `json:"hub_id" binding:"required"`
	Description string ` json:"description" binding:"max=255"`
}

type TeamSearchReq struct {
	ID   int    `form:"id"`
	Name string `form:"name"`
}

type TeamResp struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	HubID       uint    `json:"hub_id"`
	Hub         HubResp `json:"hub"`
}
