package brick

import (
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/router/auth"
)

func OptProductLotDBByAuth(db *gorm.DB, ac *auth.AccessControl) *gorm.DB {
	switch auth.RoleSet[ac.User.Role] {
	case auth.UserUserRole:
		return db.Omit("SupplierID", "Supplier", "PurchasingPrice", "Paid")
	default:
		return db
	}
}

func OptProductLotRespByAuth(resp *types.ProductLotForDetail, ac *auth.AccessControl) {
	if auth.RoleSet[ac.User.Role] == auth.UserUserRole {
		optProductLotRespForUser(resp)
	}
}

func optProductLotRespForUser(resp *types.ProductLotForDetail) {
	resp.SupplierId = ""
	resp.Supplier = nil
	resp.PurchasingPrice = 0
	resp.Paid = 0
}

func OptProductDBByAuth(db *gorm.DB, ac *auth.AccessControl) *gorm.DB {
	switch auth.RoleSet[ac.User.Role] {
	case auth.UserUserRole:
		return db
	default:
		return db.Preload("ProductLots.Supplier")
	}
}

// OptProductLotListRespByAuth
// For those functions which contain the productLot list in response
func OptProductLotListRespByAuth(resp []types.ProductLotForDetail, ac *auth.AccessControl) {
	if auth.RoleSet[ac.User.Role] == auth.UserUserRole {
		for k := range resp {
			optProductLotRespForUser(&resp[k])
		}
	}
}
