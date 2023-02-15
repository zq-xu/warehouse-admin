package supplier

import (
	"github.com/gin-gonic/gin"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/restapi"
)

func DeleteSupplier(ctx *gin.Context) {
	restapi.Delete(ctx, &model.Deliverer{})
}
