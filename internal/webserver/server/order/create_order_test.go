package order

import (
	"encoding/json"
	"testing"
)

// go test -v create_order_test.go create_order.go -test.run TestCreateOrderReq
func TestCreateOrderReq(t *testing.T) {
	r := &CreateOrderReq{}

	b, _ := json.Marshal(r)
	t.Logf("config is%+v", string(b))
}
