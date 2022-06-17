package sdk

import (
	"context"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/sirupsen/logrus"
	"time"
)

type ServerCenterHandler struct {
}

func (this *ServerCenterHandler) ListAddress(ctx context.Context) []string {
	return this.address
}
func (this *ServerCenterHandler) GetSecret(ctx context.Context) string {
	return Config.Secret
}
func (this *ServerCenterHandler) GetServerName(ctx context.Context) string {
	return sdk.GetEnvServerName(ctx, model.DefaultServerName)
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
	logrus.WithContext(ctx).WithFields(logrus.Fields{"Config": Config}).Info("加载配置")
	return nil
}
func (this *ServerCenterHandler) GetDefaultConf(ctx context.Context) string {
	var config model.Config
	config, _ = checkAndResetConfig(ctx, config)
	return util.ToYamlString(config)
}
func (this *ServerCenterHandler) GetLocalFilePath(ctx context.Context) string {
	return ""
}
