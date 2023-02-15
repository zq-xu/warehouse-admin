package log

import (
	"os"

	"github.com/sirupsen/logrus"

	"zq-xu/warehouse-admin/pkg/utils"
)

var Logger = logrus.New()

func InitLogger() {
	logrus.SetOutput(os.Stdout)
	Logger.SetLevel(LogrusCfg.Level)
	utils.Logger.Info("Succeed to init log!")
}
