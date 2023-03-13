package config

import (
	"os"

	"github.com/spf13/pflag"

	"zq-xu/warehouse-admin/pkg/config"
	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/utils"
)

const (
	webServerConfigName = "WebServerConfig"

	tmpDirEnv          = "TmpDir"
	thumbnailWidthEnv  = "ThumbnailWidth"
	thumbnailHeightEnv = "ThumbnailHeight"

	defaultTmpDir = "/webserver-tmp"
)

type WebServerConfig struct {
	TmpDir string

	ThumbnailWidth  int
	ThumbnailHeight int
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
	fs.IntVar(&wsc.ThumbnailWidth, "thumbnail-width", utils.GetIntFromEnv(thumbnailWidthEnv), "the width of the thumbnail.")
	fs.IntVar(&wsc.ThumbnailHeight, "thumbnail-height", utils.GetIntFromEnv(thumbnailHeightEnv), "the height of the thumbnail.")
}

func (wsc *WebServerConfig) Revise() {
	if wsc.TmpDir != "" {
		err := utils.EnsureDirExist(wsc.TmpDir)
		if err != nil {
			panic(err)
		}
	}

	log.Logger.Infof("WebServerConfig is %+v", wsc)
}
