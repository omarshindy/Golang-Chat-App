package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Application struct {
	Base
	Token      string `gorm:"column:token"`
	Name       string `gorm:"column:name" json:"-"`
	ChatsCount int32  `gorm:"column:chats_count" json:"-"`
}

func (application *Application) BeforeCreate(tx *gorm.DB) (err error) {
	application.Token = uuid.NewString()
	return
}
