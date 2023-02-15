package store

import (
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/pkg/utils"
)

var (
	modelSet = make([]interface{}, 0)
)

func RegisterModel(m interface{}) {
	if utils.IsInterfaceValueNil(m) {
		return
	}

	modelSet = append(modelSet, m)
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(modelSet...)
}
