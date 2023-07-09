package model

import (
	"time"

	"gorm.io/gorm"

	"zq-xu/warehouse-admin/pkg/store"
)

const (
	OrderTableName = "order"
)

var (
	OrderOmitCreate = []string{"Salesman", "Deliverer", "Customer"}
)

// description:"default the column ID is the primaryKey
type Order struct {
	store.Model

	// The Order belongs to the Customer
	CustomerID int64
	Customer   Customer

	// The Order belongs to the Salesman
	SalesmanID int64 `gorm:"default:null"`
	Salesman   Salesman

	// The Order belongs to the Deliverer
	DelivererID int64 `gorm:"default:null"`
	Deliverer   Deliverer

	OrderProducts []OrderProduct

	OrderNo         string
	Phone           string
	PayMode         int
	DeliveryMode    int
	DeliveryAddress string
	DeliveryAt      *time.Time
	Paid            float32
	Comment         string `gorm:"size:512"`
	Status          int
}

type OrderDetail struct {
	Order

	TotalPrice float32
	TotalPaid  float32
}

func (o *Order) TableName() string {
	return OrderTableName
}

func init() {
	store.RegisterTable(&Order{})
}

func GenerateReadOrderDB(db, query *gorm.DB) *gorm.DB {
	return db.
		Preload("Customer").
		Preload("Salesman").
		Preload("Deliverer").
		Preload("OrderProducts").
		Preload("OrderProducts.Product").
		Select("`order`.*,q.total_paid,q.total_price").
		Joins("left join (?) q on q.order_id = order.id", query)
}

func GenerateOrderAssociationsQuery(db *gorm.DB) *gorm.DB {
	return db.Table(OrderProductTableName).
		Select("order_id," +
			"sum(order_product.paid) as total_paid," +
			"sum(order_product.final_price) as total_price").
		Group("order_product.order_id")
}

func LoadOrderListAssociations(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Orders").
		Preload("Orders.Customer").
		Preload("Orders.Salesman").
		Preload("Orders.Deliverer").
		Preload("Orders.OrderProducts").
		Preload("Orders.OrderProducts.Product")
}
