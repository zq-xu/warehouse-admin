package product

import (
	"github.com/gin-gonic/gin"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/restapi"
)

func DeleteProduct(ctx *gin.Context) {
	restapi.Delete(ctx, &model.Product{})
}
