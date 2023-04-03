package salesman

import (
	"encoding/json"
	"testing"
)

// go test -v create_salesman_test.go create_salesman.go -test.run TestCreateSalesmanReq
func TestCreateSalesmanReq(t *testing.T) {
	r := &CreateSalesmanReq{}

	b, _ := json.Marshal(r)
	t.Logf("config is%+v", string(b))
}
