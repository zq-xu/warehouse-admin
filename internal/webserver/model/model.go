package model

import (
	"time"

	"gorm.io/gorm"

	"zq-xu/warehouse-admin/pkg/utils"
)

type Model struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func GenerateModel() Model {
	return Model{
		ID:        utils.GenerateUUID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
