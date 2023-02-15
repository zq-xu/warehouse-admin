package main

import (
	"fmt"
	"zq-xu/warehouse-admin/pkg/store"

	"zq-xu/warehouse-admin/internal/webserver/config"
	_ "zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/internal/webserver/router"
	"zq-xu/warehouse-admin/pkg/log"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic error is %v", err)
		}
	}()

	config.InitConfig()

	log.InitLogger()

	store.InitDatabase()

	assert(router.StartRouter())
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}
