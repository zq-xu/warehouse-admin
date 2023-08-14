package order

import (
	"github.com/gin-gonic/gin"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/restapi"
)

func DeleteOrder(ctx *gin.Context) {
	restapi.Delete(ctx, &model.Deliverer{})
}
