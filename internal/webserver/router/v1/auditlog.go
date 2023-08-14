package v1

import (
	"net/http"

	"zq-xu/warehouse-admin/pkg/router"
	"zq-xu/warehouse-admin/pkg/router/auditlog"
)

var (
	auditLogGroupPath = "/auditlogs"
	AuditLogGroup     = &router.APIGroup{
		RelativePath: auditLogGroupPath,
		APIs: []*router.API{
			{Method: http.MethodGet, Path: "", Handler: auditlog.ListAuditLog},
			{Method: http.MethodGet, Path: "/export", Handler: auditlog.ExportAuditLog},
		},
	}
)

func init() {
	registerAPIGroup(AuditLogGroup)
}
