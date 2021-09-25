package config

import (
	"context"
	"fmt"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/cellargalaxy/server_center/sdk"
	"github.com/sirupsen/logrus"
	"time"
)

var Config = model.Config{}

func init() {
	var err error
	ctx := util.CreateLogCtx()
	Config.LogLevel = logrus.InfoLevel
	Config.MysqlDsn = util.GetEnvString("mysql_dsn", Config.MysqlDsn)
	Config.ShowSql = false
	Config.Secret = "secret"
	Config, err = checkAndResetConfig(ctx, Config)
	if err != nil {
		panic(err)
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{"Config": Config}).Info("加载配置")
}

func Init(ctx context.Context) {
	client, _ := sdk.NewDefaultServerCenterClient(ctx, &ServerCenterHandler{})
	if client == nil {
		return
	}
	client.StartConf(ctx)
}

func checkAndResetConfig(ctx context.Context, config model.Config) (model.Config, error) {
	if config.LogLevel <= 0 || config.LogLevel > logrus.TraceLevel {
		config.LogLevel = logrus.InfoLevel
	}
	if config.MysqlDsn == "" {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Error("mysql_dsn为空")
		return config, fmt.Errorf("mysql_dsn为空")
	}
	if config.Secret == "" {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Error("secret为空")
		return config, fmt.Errorf("secret为空")
	}
	return config, nil
}

type ServerCenterHandler struct {
}

func (this *ServerCenterHandler) GetAddress(ctx context.Context) string {
	return "http://127.0.0.1" + model.ListenAddress
}
func (this *ServerCenterHandler) GetSecret(ctx context.Context) string {
	return Config.Secret
}
func (this *ServerCenterHandler) GetServerName(ctx context.Context) string {
	return sdk.GetEnvServerName(ctx)
}
func (this *ServerCenterHandler) GetInterval(ctx context.Context) time.Duration {
	return 5 * time.Minute
}
func (this *ServerCenterHandler) ParseConf(ctx context.Context, object model.ServerConfModel) error {
	var config model.Config
	err := util.UnmarshalYamlString(object.ConfText, &config)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("反序列化配置异常")
		return err
	}
	config, err = checkAndResetConfig(ctx, config)
	if err != nil {
		return err
	}
	Config = config
	logrus.SetLevel(Config.LogLevel)
	logrus.WithContext(ctx).WithFields(logrus.Fields{"Config": Config}).Info("加载配置")
	return nil
}
