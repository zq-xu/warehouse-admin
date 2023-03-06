package model

import (
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/pkg/store"
)

const (
	CustomerTableName = "customer"
)

// description:"default the column ID is the primaryKey
type Customer struct {
	store.Model

	Name        string
	Phone       string
	Address     string
	InvoiceInfo string
	Comment     string
	Status      int

	Orders []Order
}

type CustomerDetail struct {
	Customer

	TotalPrice float32
	TotalPaid  float32
}

func (o *Customer) TableName() string {
	return CustomerTableName
}

func init() {
	store.RegisterTable(&Customer{})
}

func GenerateReadCustomerDB(db *gorm.DB, query *gorm.DB) *gorm.DB {
	return db.
		Select("customer.*,q.total_paid,q.total_price").
		Joins("left join (?) q on q.customer_id = customer.id", query)
}

func GenerateCustomerAssociationsQuery(db *gorm.DB) *gorm.DB {
	opQuery := db.Table(OrderProductTableName).
		Select("order_id," +
			"sum(order_product.paid) as total_paid," +
			"sum(order_product.bought_price) as total_price").
		Group("order_product.order_id")

	return db.Table(OrderTableName).
		Select("customer_id,"+
			"sum(q2.total_paid) as total_paid,"+
			"sum(q2.total_price) as total_price").
		Group("order.customer_id").
		Joins("left join (?) q2 on q2.order_id = order.id", opQuery)
}
