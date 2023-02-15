package model

import (
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/pkg/store"
)

const (
	SalesmanTableName = "salesman"
)

// description:"default the column ID is the primaryKey
type Salesman struct {
	Model

	Name    string
	Phone   string
	Comment string
	Status  int

	Orders []Order
}

type SalesmanDetail struct {
	Salesman

	TotalPrice float32
	TotalPaid  float32
}

func (o *Salesman) TableName() string {
	return SalesmanTableName
}

func init() {
	store.RegisterModel(&Salesman{})
}

func GenerateReadSalesmanDB(db *gorm.DB, query *gorm.DB) *gorm.DB {
	return db.
		Select("salesman.*,q.total_paid,q.total_price").
		Joins("left join (?) q on q.salesman_id = salesman.id", query)
}

func GenerateSalesmanAssociationsQuery(db *gorm.DB) *gorm.DB {
	opQuery := db.Table(OrderProductTableName).
		Select("order_id," +
			"sum(order_product.paid) as total_paid," +
			"sum(order_product.current_price) as total_price").
		Group("order_product.order_id")

	return db.Table(OrderTableName).
		Select("salesman_id,"+
			"sum(q2.total_paid) as total_paid,"+
			"sum(q2.total_price) as total_price").
		Group("order.salesman_id").
		Joins("left join (?) q2 on q2.order_id = order.id", opQuery)
}
