package config

import (
	"zq-xu/warehouse-admin/pkg/config"
	"zq-xu/warehouse-admin/pkg/store"
)

func InitConfig() {
	store.InitDataBaseConfig()

	InitWebServerConfig()

	config.InitConfig()
}
