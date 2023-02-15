package router

import (
	"github.com/spf13/pflag"

	"zq-xu/warehouse-admin/pkg/config"
)

const (
	routerConfigName = "RouterConfig"
)

type RouteConfig struct {
	IP        string
	Port      string
	PprofPort string
}

var (
	RouteCfg = &RouteConfig{}
)

func init() {
	config.RegisterCfg(routerConfigName, RouteCfg)
}

// AddFlags adds flags for router
func (rc *RouteConfig) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&rc.IP, "router-ip", "0.0.0.0", "the ip for listening")
	fs.StringVar(&rc.Port, "router-port", "8080", "the port for listening")
	fs.StringVar(&rc.PprofPort, "pprof-port", "6069", "the port for pprof")
}

func (rc *RouteConfig) Revise() {}
