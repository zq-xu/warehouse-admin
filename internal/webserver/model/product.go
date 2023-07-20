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
	store.Model

	Name           string
	Image          string
	Thumbnail      string
	Price          float32
	StorageAddress string
	Comment        string `gorm:"size:512"`
	Status         int

	// The Order belongs to the Deliverer
	CategoryID *int64 `gorm:"default:null"`
	Category   Category

	// has many relation: used to preload the ProductLots
	ProductLots []ProductLot

	// has many relation: used to preload the OrderProducts
	OrderProducts []OrderProduct
}

type ProductDetail struct {
	Product

	TotalCount int
	SoldCount  int
	Stocks     int
}

func (o *Product) TableName() string {
	return ProductTableName
}

func init() {
	store.RegisterTable(&Product{})
}

func GenerateReadProductDB(db, queryDB *gorm.DB) *gorm.DB {
	return queryDB.
		Select("product.*,q.total_count,op.sold_count,(q.total_count - op.sold_count) as stocks").
		Joins("left join (?) q on q.product_id = product.id", generateProductTotalCountQuery(db)).
		Joins("left join (?) op on op.product_id = product.id", generateProductSoldCountQuery(db))
}

func generateProductTotalCountQuery(db *gorm.DB) *gorm.DB {
	return db.Table(ProductLotTableName).
		Select("product_id," +
			"sum(product_lot.count) as total_count").
		Group("product_lot.product_id")
}

func generateProductSoldCountQuery(db *gorm.DB) *gorm.DB {
	return db.Table(OrderProductTableName).
		Select("product_id," +
			"sum(order_product.count) as sold_count").
		Group("order_product.product_id")
}
