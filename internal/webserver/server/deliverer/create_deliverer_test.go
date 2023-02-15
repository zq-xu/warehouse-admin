package deliverer

import (
	"encoding/json"
	"testing"
)

//  go test -v create_deliverer_test.go create_deliverer.go -test.run TestCreateDelivererReq
func TestCreateDelivererReq(t *testing.T) {
	r := &CreateDelivererReq{}

	b, _ := json.Marshal(r)
	t.Logf("config is%+v", string(b))
}
