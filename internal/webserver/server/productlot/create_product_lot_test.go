package productlot

import (
	"encoding/json"
	"testing"
)

//  go test -v create_product_lot_test.go create_product_lot.go -test.run TestCreateProductLotReq
func TestCreateProductLotReq(t *testing.T) {
	r := &CreateProductLotReq{}

	b, _ := json.Marshal(r)
	t.Logf("config is%+v", string(b))
}
