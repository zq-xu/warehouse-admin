package category

import (
	"github.com/gin-gonic/gin"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/restapi"
)

func DeleteCategory(ctx *gin.Context) {
	restapi.Delete(ctx, &model.Category{})
}
