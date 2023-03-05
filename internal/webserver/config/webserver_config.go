package config

import (
	"os"

	"github.com/spf13/pflag"

	"zq-xu/warehouse-admin/pkg/config"
	"zq-xu/warehouse-admin/pkg/log"
)

const (
	webServerConfigName = "WebServerConfig"

	tmpDirEnv     = "TmpDir"
	defaultTmpDir = "/webserver-tmp"
)

type WebServerConfig struct {
	NodePortIP string
	TmpDir     string
}

var (
	WebServerCfg = &WebServerConfig{}
)

func InitWebServerConfig() {
	config.RegisterCfg(webServerConfigName, WebServerCfg)
}

// AddFlags adds flags for router
func (wsc *WebServerConfig) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&wsc.TmpDir, "tmp-dir", os.Getenv(tmpDirEnv), "the tmp dir to store temporary files.")
}

func (wsc *WebServerConfig) Revise() {
	if wsc.TmpDir == "" {
		wsc.TmpDir = defaultTmpDir
	}

	log.Logger.Infof("WebServerConfig is %+v", wsc)
}
