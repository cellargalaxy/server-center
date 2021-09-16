package config

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/sirupsen/logrus"
)

var Config = model.Config{}

func init() {
	ctx := util.CreateLogCtx()
	Config.LogLevel = logrus.Level(util.GetEnvInt("log_level", int(Config.LogLevel)))
	Config.MysqlDsn = util.GetEnvString("mysql_dsn", Config.MysqlDsn)
	Config.ShowSql = util.GetEnvBool("show_sql", Config.ShowSql)
	Config.Secret = util.GetEnvString("secret", Config.Secret)
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
