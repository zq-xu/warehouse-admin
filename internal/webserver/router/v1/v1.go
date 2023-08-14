package v1

import (
	"zq-xu/warehouse-admin/pkg/router"
	"zq-xu/warehouse-admin/pkg/router/auditlog"
	"zq-xu/warehouse-admin/pkg/router/auth"
	"zq-xu/warehouse-admin/pkg/utils"
)

const (
	VersionV1 = "/v1"
)

var apiGrps = make([]*router.APIGroup, 0)

func registerAPIGroup(grps ...*router.APIGroup) {
	for _, grp := range grps {
		if utils.IsInterfaceValueNil(grp) {
			return
		}

		apiGrps = append(apiGrps, grp)
	}
}

func Register() {
	v1Grp := router.NewGroup(VersionV1)
	v1Grp.AddMiddlewares(auth.Middleware.MiddlewareFunc())
	v1Grp.AddMiddlewares(auditlog.Middleware.MiddlewareFunc())

	for _, v := range apiGrps {
		v1Grp.AddAPIGroup(v)
	}

	router.RegisterGroup(v1Grp)
}
