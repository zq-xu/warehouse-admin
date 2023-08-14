package product

import (
	"encoding/json"
	"testing"
)

// go test -v create_product_test.go create_product.go -test.run TestCreateProductReq
func TestCreateProductReq(t *testing.T) {
	r := &CreateProductReq{}

	b, _ := json.Marshal(r)
	t.Logf("config is%+v", string(b))
}
