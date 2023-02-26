package model

import (
	"time"

	"zq-xu/warehouse-admin/pkg/store"
)

const (
	ProductLotTableName = "product_lot"
)

// description:"default the column ID is the primaryKey
type ProductLot struct {
	Model

	// The ProductLot belongs to the Product
	ProductID int64
	Product   Product

	// The ProductLot belongs to the Supplier
	SupplierID int64
	Supplier   Supplier

	LotNo           string
	PurchasingPrice float32
	PurchasingDate  *time.Time
	Paid            float32
	Count           int
	StorageAddress  string
	Comment         string `gorm:"size:512"`
}

func (o *ProductLot) TableName() string {
	return ProductLotTableName
}

func init() {
	store.RegisterModel(&ProductLot{})
}
