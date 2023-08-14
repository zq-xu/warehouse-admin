package model

import (
	"zq-xu/warehouse-admin/pkg/store"
)

const (
	StockOutTableName = "stock_out"
)

// description:"default the column ID is the primaryKey
type StockOut struct {
	store.Model

	// The StockOut belongs to the Order
	OrderID int64
	Order   Order

	// The StockOut belongs to the ProductLot
	ProductLotID int64
	ProductLot   ProductLot

	Count int

	Comment string `gorm:"size:512"`
}

func (op *StockOut) TableName() string {
	return StockOutTableName
}

func init() {
	store.RegisterTable(&StockOut{})
}
