package customer

import (
	"encoding/json"
	"testing"
)

//  go test -v create_customer_test.go create_customer.go -test.run TestCreateCustomerReq
func TestCreateCustomerReq(t *testing.T) {
	r := &CreateCustomerReq{}

	b, _ := json.Marshal(r)
	t.Logf("config is%+v", string(b))
}
