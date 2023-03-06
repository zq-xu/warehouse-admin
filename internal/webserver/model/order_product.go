package model

import (
	"zq-xu/warehouse-admin/pkg/store"
)

const (
	OrderProductTableName = "order_product"
)

// description:"default the column ID is the primaryKey
type OrderProduct struct {
	store.Model

	// The OrderProduct belongs to the Order
	OrderID int64
	Order   Order

	// The OrderProduct belongs to the Product
	ProductID int64
	Product   Product

	Count int
	// The price of the product is changeable,
	// So it's necessary to record it.
	BoughtPrice float32
	Paid        float32
	Discount    int
	Comment     string `gorm:"size:512"`
}

func (op *OrderProduct) TableName() string {
	return OrderProductTableName
}

func init() {
	store.RegisterTable(&OrderProduct{})
}

// For Copier
func (op *OrderProduct) Name() string {
	return op.Product.Name
}

func (op *OrderProduct) Price() float32 {
	return op.Product.Price
}
