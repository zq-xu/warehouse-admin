package model

import (
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/pkg/store"
)

const (
	CategoryTableName = "category"
)

var (
	CategoryBaseInfoColumns = []string{"id", "name"}
)

// description:"default the column ID is the primaryKey
type Category struct {
	store.Model

	Name string

	Products []ProductDetail
}

type CategoryDetail struct {
	Category

	ProductCount int
}

func (o *Category) TableName() string {
	return CategoryTableName
}

func init() {
	store.RegisterTable(&Category{})
}

func GenerateReadCategoryDetailDB(db *gorm.DB) *gorm.DB {
	return categoryLoadProductCount(db, categoryLoadProductsList(db))
}

func categoryLoadProductsList(db *gorm.DB) *gorm.DB {
	return db.Preload("Products", func(innerDB *gorm.DB) *gorm.DB {
		return GenerateReadProductDB(db, innerDB)
	})
}

func GenerateReadCategoryListDB(db, queryDB *gorm.DB) *gorm.DB {
	return categoryLoadProductCount(db, queryDB)
}

func categoryLoadProductCount(db, queryDB *gorm.DB) *gorm.DB {
	return queryDB.
		Select("category.*,q.product_count").
		Joins("left join (?) q on q.category_id = category.id", generateCategoryAssociationsQuery(db))
}

func generateCategoryAssociationsQuery(db *gorm.DB) *gorm.DB {
	opQuery := db.Table(ProductLotTableName).
		Select("product_id," +
			"sum(product_lot.count) as product_count").
		Group("product_lot.product_id")

	return db.Table(ProductTableName).
		Select("category_id,"+
			"sum(q2.product_count) as product_count").
		Group("product.category_id").
		Joins("left join (?) q2 on q2.product_id = product.id", opQuery)
}
