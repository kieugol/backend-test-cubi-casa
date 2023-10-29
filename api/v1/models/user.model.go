package models

import (
	"time"

	"github.com/backend-test-cubi-casa/helpers/util"
)

func (m *User) GetAlias() string {
	return util.GetStructName(User{})
}

type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Phone     string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone"`
	Address   string    `gorm:"type:varchar(25);not null" json:"address"`
	TeamID    uint      `json:"team_id"`
	Team      Team      `gorm:"foreignKey:TeamID" json:"-"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

type UserCreateReq struct {
	Name    string `json:"name" binding:"required,max=255"`
	Email   string `json:"email" binding:"required,email"`
	TeamID  uint   `json:"team_id" binding:"required"`
	Phone   string `json:"phone" binding:"required,e164"`
	Address string `json:"address" binding:"required,max=255"`
}

type UserSearchReq struct {
	Name  string `form:"name"`
	Email string `form:"email"`
}

type UserResp struct {
	ID      uint     `json:"id"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Phone   string   `json:"phone"`
	Address string   `json:"address"`
	TeamID  uint     `json:"team_id"`
	Team    TeamResp `json:"team"`
}
