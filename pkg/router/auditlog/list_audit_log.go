package auditlog

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/router/auth"
)

type ResponseOfAuditLog struct {
	types.ModelBase `json:",inline"`

	User auth.ResponseOfUserInfo `json:"user"`

	ClientIP   string `json:"client_ip"`
	Url        string `json:"url"`
	Method     string `json:"method"`
	Body       string `json:"body"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func ListAuditLog(ctx *gin.Context) {
	listObj := make([]ModelAuditLog, 0)

	conf := &restapi.ListConf{
		ModelObj:              &ModelAuditLog{},
		ModelObjList:          &listObj,
		FuzzySearchColumnList: []string{"url"},
		AuthControl: restapi.AuthControl{
			AuthValidation: func(ac *auth.AccessControl) bool { return ac.User.Role > 0 },
		},
		TransObjToRespFunc:     func(ac *auth.AccessControl) []interface{} { return generateAuditLogListResponse(listObj) },
		LoadAssociationsDBFunc: loadAssociationsDB,
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func loadAssociationsDB(db, queryDB *gorm.DB, ac *auth.AccessControl) *gorm.DB {
	return GenerateReadAuditLogDB(queryDB)
}

func generateAuditLogListResponse(objList []ModelAuditLog) []interface{} {
	items := make([]interface{}, 0)

	for _, v := range objList {
		r := ResponseOfAuditLog{}

		err := copier.Copy(&r, &v)
		if err != nil {
			log.Logger.Errorf("Failed to copy audit log obj to response. %v", err)
			return nil
		}

		items = append(items, r)
	}

	return items
}
