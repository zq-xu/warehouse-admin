package router

import (
	"net"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"zq-xu/warehouse-admin/pkg/utils"
)

var groups = make([]*APIGroup, 0)

type engineOpt func(r *gin.Engine)

func RegisterGroup(grps ...*APIGroup) {
	for _, grp := range grps {
		if utils.IsInterfaceValueNil(grp) {
			return
		}

		groups = append(groups, grp)
	}
}

func StartRouter(r *gin.Engine) error {
	for _, grp := range groups {
		grp.AddToEngine(r)
	}

	return r.Run(net.JoinHostPort(RouteCfg.IP, RouteCfg.Port))
}

func NewRouter(fns ...engineOpt) *gin.Engine {
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	for _, fn := range fns {
		fn(r)
	}

	return r
}

func DefaultRouter() *gin.Engine {
	return NewRouter(func(r *gin.Engine) {
		r.Use(LoggerFilter([]string{HealthPath}, GetMethodFilter))
		r.Use(gzip.Gzip(gzip.DefaultCompression))
		r.Use(corsMiddleware())
		registerHealth(r)
	})
}
