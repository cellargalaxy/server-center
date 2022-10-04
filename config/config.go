package config

import (
	"context"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/cellargalaxy/server_center/sdk"
	"github.com/sirupsen/logrus"
	"time"
)

var Config = model.Config{}

func init() {
	ctx := util.GenCtx()
	var err error

	var handler ServerCenterHandler
	client, err := sdk.NewDefaultServerCenterClient(ctx, &handler)
	if err != nil {
		panic(err)
	}
	if client == nil {
		panic("创建ServerCenterClient为空")
	}
	err = client.StartWithInitConf(ctx)
	if err != nil {
		panic(err)
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{"Config": Config}).Info("加载配置")
}

func checkAndResetConfig(ctx context.Context, config model.Config) (model.Config, error) {
	if config.ClearEventSave <= 0 {
		config.ClearEventSave = 1000000
	}
	if config.ClearConfigSave <= 0 {
		config.ClearConfigSave = 100
	}

	if config.Secret == "" {
		go func() {
			for i := 0; i < 100; i++ {
				logrus.WithContext(ctx).WithFields(logrus.Fields{}).Warn("加载配置，Secret为空")
				util.Sleep(ctx, time.Second*10)
			}
		}()
	}

	return config, nil
}

type ServerCenterHandler struct {
	sdk.ServerCenterDefaultHandler
}

func (this *ServerCenterHandler) GetServerName(ctx context.Context) string {
	return sdk.GetEnvServerName(ctx, model.DefaultServerName)
}
func (this *ServerCenterHandler) ParseConf(ctx context.Context, object model.ServerConfModel) error {
	var config model.Config
	err := util.UnmarshalYamlString(object.ConfText, &config)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"ServerName": object.ServerName, "Version": object.Version, "err": err}).Error("反序列化配置异常")
		return err
	}
	config, err = checkAndResetConfig(ctx, config)
	if err != nil {
		return err
	}
	Config = config
	logrus.WithContext(ctx).WithFields(logrus.Fields{"Config": Config}).Info("加载配置")
	return nil
}
func (this *ServerCenterHandler) GetDefaultConf(ctx context.Context) string {
	var config model.Config
	config, _ = checkAndResetConfig(ctx, config)
	return util.ToYamlString(config)
}
