package v1

import (
	"zq-xu/warehouse-admin/pkg/router"
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
	//v1Grp := router.NewGroup(VersionV1).
	//	AddMiddlewares(auth.AuthMiddleware.MiddlewareFunc())

	v1Grp := router.NewGroup(VersionV1)

	for _, v := range apiGrps {
		v1Grp.AddAPIGroup(v)
	}

	router.RegisterGroup(v1Grp)
}