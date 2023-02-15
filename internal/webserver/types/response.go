package types

import (
	"fmt"
	"time"
	"zq-xu/warehouse-admin/pkg/utils"
)

type ModelBase struct {
	Id        string         `json:"id"`
	CreatedAt utils.UnixTime `json:"createdAt"`
	UpdatedAt utils.UnixTime `json:"updatedAt"`
}

func (d *ModelBase) ID(id int64) {
	d.Id = fmt.Sprintf("%d", id)
}

type CustomerForDetail struct {
	ModelBase `json:",inline"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
}

type SupplierForDetail struct {
	ModelBase `json:",inline"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
}

type SalesmanForDetail struct {
	ModelBase `json:",inline"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
}

type DelivererForDetail struct {
	ModelBase `json:",inline"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
}

type ProductForDetail struct {
	ModelBase      `json:",inline"`
	Name           string  `json:"name"`
	Image          string  `json:"image"`
	Price          float32 `json:"price"`
	StorageAddress string  `json:"storageAddress"`
	Comment        string  `json:"comment"`
}

type ProductLotForDetail struct {
	ModelBase

	LotNo           string    `json:"lotNo"`
	PurchasingCount int       `json:"purchasingCount"`
	PurchasingPrice float32   `json:"purchasingPrice"`
	PurchasingDate  time.Time `json:"purchasingDate"`
	Paid            float32   `json:"paid"`
	Count           int       `json:"count"`
	StorageAddress  string    `json:"storageAddress"`
	Comment         string    `json:"comment"`
	Status          int       `json:"status"`

	Product  ProductForDetail  `json:"product"`
	Supplier SupplierForDetail `json:"supplier"`
}

type OrderForDetail struct {
	ModelBase

	CustomerId string `json:"customerId"`
	SalesmanId string `json:"salesmanId"`

	PayMethod int     `json:"payMethod"`
	Paid      float32 `json:"paid"`
	Comment   string  `json:"comment"`

	DeliveryMode    int             `json:"deliveryMode"`
	DeliveryAt      *utils.UnixTime `json:"deliveryAt"`
	DeliveryAddress string          `json:"deliveryAddress"`
	DelivererId     string          `json:"delivererId"`

	OrderNo    string  `json:"orderNo"`
	TotalPrice float32 `json:"totalPrice"`
	TotalPaid  float32 `json:"totalPaid"`
	Status     int     `json:"status"`

	Customer      CustomerForDetail       `json:"customer"`
	Salesman      SalesmanForDetail       `json:"salesman"`
	Deliverer     DelivererForDetail      `json:"deliverer"`
	OrderProducts []OrderProductForDetail `json:"products"`
}

// if order is preload, the 'TotalPrice' and 'TotalPaid' can't be calculated,
// so revise these property manually.
func (od *OrderForDetail) Revise() {
	od.TotalPrice = 0
	od.TotalPaid = 0

	for _, v := range od.OrderProducts {
		od.TotalPrice += v.Price
		od.TotalPaid += v.Paid
	}
}

type OrderProductForDetail struct {
	ModelBase   `json:",inline"`
	ProductId   string  `json:"productId"`
	Name        string  `json:"name"`
	Price       float32 `json:"price"`
	Count       int     `json:"count"`
	BoughtPrice float32 `json:"boughtPrice"`
	Paid        float32 `json:"paid"`
	Discount    int     `json:"discount"`
	Comment     string  `json:"comment"`
}

func (opd *OrderProductForDetail) ProductID(id int64) {
	opd.ProductId = fmt.Sprintf("%d", id)
}
