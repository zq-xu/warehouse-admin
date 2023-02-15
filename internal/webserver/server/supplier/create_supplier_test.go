package supplier

import (
	"encoding/json"
	"testing"
)

//  go test -v create_supplier_test.go create_supplier.go -test.run TestCreateSupplierReq
func TestCreateSupplierReq(t *testing.T) {
	r := &CreateSupplierReq{}

	b, _ := json.Marshal(r)
	t.Logf("config is%+v", string(b))
}
