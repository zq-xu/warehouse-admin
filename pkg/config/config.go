package config

import (
	"flag"
	"os"

	"github.com/spf13/pflag"

	"zq-xu/warehouse-admin/pkg/utils"
)

type cfg interface {
	AddFlags(fs *pflag.FlagSet)
	Revise()
}

var (
	cfgSet = make(map[string]cfg)
)

func RegisterCfg(n string, c cfg) {
	_, ok := cfgSet[n]
	if ok {
		//log.Loggercicd.Warningf("cfg %s has already exist", n)
		return
	}

	cfgSet[n] = c
}

func InitConfig() {
	InitFlag(1)
}

func InitConfigWithSubCommand(subCommandCount int) []string {
	subCommands := GetSubCommand(subCommandCount)
	InitFlag(subCommandCount + 1)
	return subCommands
}

func InitConfigWithSingleSubCommand() string {
	subCommands := GetSubCommand(1)
	InitFlag(2)

	if len(subCommands) < 1 {
		return ""
	}

	return subCommands[0]
}

func GetSingleCommand() string {
	subCommands := GetSubCommand(1)
	if len(subCommands) < 1 {
		return ""
	}

	return subCommands[0]
}

func GetSubCommand(count int) []string {
	if len(os.Args) < count {
		return nil
	}

	subCommands := os.Args[1 : count+1]
	utils.Logger.Infof("Sub commands are %+v", subCommands)

	return subCommands
}

func InitFlag(flagStartIndex int) {
	for _, v := range cfgSet {
		v.AddFlags(pflag.CommandLine)
	}

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	_ = pflag.CommandLine.Parse(os.Args[flagStartIndex:])

	for k, v := range cfgSet {
		v.Revise()
		utils.Logger.Infof("Config %s is %+v", k, v)
	}

	utils.Logger.Info("Succeed to init the config!")
}
