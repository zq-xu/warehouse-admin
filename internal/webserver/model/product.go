package model

import (
	"gorm.io/gorm"
	"zq-xu/warehouse-admin/pkg/store"
)

const (
	ProductTableName = "product"
)

// description:"default the column ID is the primaryKey
type Product struct {
	Model

	Name           string
	Image          string
	Price          float32
	StorageAddress string
	Comment        string `gorm:"size:512"`
	Status         int

	// has many relation: used to preload the ProductLots
	ProductLots []ProductLot

	// has many relation: used to preload the OrderProducts
	OrderProducts []OrderProduct
}

type ProductDetail struct {
	Product

	TotalCount int
}

func (o *Product) TableName() string {
	return ProductTableName
}

func init() {
	store.RegisterModel(&Product{})
}

func GenerateReadProductDB(db *gorm.DB, query *gorm.DB) *gorm.DB {
	return db.Table(ProductTableName).
		Select("product.*,q.total_count").
		Joins("left join (?) q on q.product_id = product.id", query)
}

func GenerateProductAssociationsQuery(db *gorm.DB) *gorm.DB {
	return db.Table(ProductLotTableName).
		Select("product_id," +
			"sum(product_lot.count) as total_count").
		Group("product_lot.product_id")
}
