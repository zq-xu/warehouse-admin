package category

import (
	"encoding/json"
	"testing"
)

// go test -v create_category_test.go create_category.go -test.run TestCreateCategoryReq
func TestCreateCategoryReq(t *testing.T) {
	r := &CreateCategoryReq{}

	b, _ := json.Marshal(r)
	t.Logf("config is%+v", string(b))
}
