package entities

import (
	"time"

	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	ID        int32     `gorm:"primaryKey; column:id" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
}

func (base *Base) BeforeCreate(tx *gorm.DB) (err error) {
	base.CreatedAt = time.Now()
	base.UpdatedAt = time.Now()
	return
}
