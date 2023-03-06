package model

import (
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/pkg/store"
)

const (
	DelivererTableName = "deliverer"
)

// description:"default the column ID is the primaryKey
type Deliverer struct {
	store.Model

	Name    string
	Phone   string
	Comment string
	Status  int

	Orders []Order
}

func (o *Deliverer) TableName() string {
	return DelivererTableName
}

func init() {
	store.RegisterTable(&Deliverer{})
}

func GenerateReadDelivererDB(db *gorm.DB) *gorm.DB {
	return db
}
