package config

import (
	"zq-xu/warehouse-admin/pkg/awsapi"
	"zq-xu/warehouse-admin/pkg/config"
	"zq-xu/warehouse-admin/pkg/store"
)

func InitConfig() {
	store.InitDataBaseConfig()

	awsapi.InitSessionConfig()
	awsapi.InitS3Config()

	InitWebServerConfig()
	config.InitConfig()
}
