package config

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server-center/model"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
)

const (
	configFilePath = "resources/config.ini"
)

var Config = model.Config{}

func init() {
	ctx := util.CreateLogCtx()
	exist, _ := util.ExistAndIsFile(ctx, configFilePath)
	if exist {
		cfg, err := ini.Load(configFilePath)
		if err != nil {
			panic(err)
		}
		err = cfg.MapTo(&Config)
		if err != nil {
			panic(err)
		}
	}
	checkAndResetConfig()
	logrus.WithContext(ctx).WithFields(logrus.Fields{"Config": Config}).Info("加载配置")
}

func checkAndResetConfig() {
	if Config.LogLevel <= 0 || Config.LogLevel > logrus.TraceLevel {
		Config.LogLevel = logrus.InfoLevel
	}
	if Config.MysqlDsn == "" {
		panic("MysqlDsn为空")
	}
	if Config.Secret == "" {
		panic("Secret为空")
	}
}
