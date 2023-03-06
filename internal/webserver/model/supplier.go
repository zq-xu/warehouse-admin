package model

import (
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/pkg/store"
)

const (
	SupplierTableName = "supplier"
)

// description:"default the column ID is the primaryKey
type Supplier struct {
	store.Model

	Name        string
	Phone       string
	Address     string
	InvoiceInfo string
	Comment     string
	Status      int

	ProductLots []ProductLot
}

type SupplierDetail struct {
	Supplier

	TotalPrice float32
	TotalPaid  float32
}

func (sr *Supplier) TableName() string {
	return SupplierTableName
}

func init() {
	store.RegisterTable(&Supplier{})
}

func GenerateReadSupplierDB(db *gorm.DB, query *gorm.DB) *gorm.DB {
	return db.
		Preload("ProductLots.Product").
		Select("supplier.*,q.total_paid,q.total_price").
		Joins("left join (?) q on q.supplier_id = supplier.id", query)
}

func GenerateSupplierAssociationsQuery(db *gorm.DB) *gorm.DB {
	return db.Table(ProductLotTableName).
		Select("supplier_id," +
			"sum(product_lot.paid) as total_paid," +
			"sum(product_lot.purchasing_price) as total_price").
		Group("product_lot.supplier_id")
}
