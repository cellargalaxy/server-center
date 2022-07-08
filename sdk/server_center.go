package sdk

import (
	"context"
	"fmt"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/sirupsen/logrus"
	"time"
)

type ServerCenterHandler struct {
}

func (this *ServerCenterHandler) GetAddress(ctx context.Context) string {
	return GetEnvServerCenterAddress(ctx)
}
func (this *ServerCenterHandler) GetSecret(ctx context.Context) string {
	return GetEnvServerCenterSecret(ctx)
}
func (this *ServerCenterHandler) GetServerName(ctx context.Context) string {
	return model.DefaultServerName
}
func (this *ServerCenterHandler) GetInterval(ctx context.Context) time.Duration {
	return time.Minute * 5
}
func (this *ServerCenterHandler) ParseConf(ctx context.Context, object model.ServerConfModel) error {
	var config model.Config
	err := util.UnmarshalYamlString(object.ConfText, &config)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("反序列化server_center配置异常")
		return err
	}
	secret = config.Secret

	list := addresses
	list = append(list, config.Addresses...)
	list = util.DistinctString(ctx, list)

	local := "http://127.0.0.1" + model.ListenAddress
	ls := make([]string, 0, len(list))
	for i := range list {
		if list[i] == "" || list[i] == local {
			continue
		}
		_, err = serverCenterClient.Ping(ctx, list[i])
		if err != nil {
			continue
		}
		ls = append(ls, list[i])
	}
	if this.GetServerName(ctx) == GetEnvServerName(ctx, this.GetServerName(ctx)) {
		ls = append(ls, local)
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{"ls": ls}).Info("加载server_center地址")
	if len(ls) == 0 {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Error("server_center地址为空")
		return fmt.Errorf("server_center地址为空")
	}
	addresses = ls
	return nil
}
func (this *ServerCenterHandler) GetDefaultConf(ctx context.Context) string {
	var config model.Config
	return util.ToYamlString(config)
}
