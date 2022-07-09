package sdk

import (
	"context"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/sirupsen/logrus"
	"time"
)

var serverCenterClient *ServerCenterClient
var addresses []string
var secret string

func initServerCenter(ctx context.Context) {
	var err error

	address := GetEnvServerCenterAddress(ctx)
	if address != "" {
		addresses = append(addresses, address)
	}
	secret = GetEnvServerCenterSecret(ctx)

	var handler ServerCenterHandler
	serverCenterClient, err = NewDefaultServerCenterClient(ctx, &handler)
	if err != nil {
		panic(err)
	}
	if serverCenterClient == nil {
		panic("创建serverCenterClient为空")
	}
	serverCenterClient.StartConfWithInitConf(ctx)
}

func ListAddress(ctx context.Context) []string {
	return addresses
}
func GetSecret(ctx context.Context) string {
	return secret
}

type ServerCenterDefaultHandler struct {
	intervalIndex int
}

func (this *ServerCenterDefaultHandler) ListAddress(ctx context.Context) []string {
	return ListAddress(ctx)
}
func (this *ServerCenterDefaultHandler) GetSecret(ctx context.Context) string {
	return GetSecret(ctx)
}
func (this *ServerCenterDefaultHandler) GetInterval(ctx context.Context) time.Duration {
	intervals := []time.Duration{time.Second * 10, time.Second * 10, time.Second * 10, time.Minute * 10}
	index := this.intervalIndex % len(intervals)
	this.intervalIndex = this.intervalIndex + 1
	return intervals[index]
}

type ServerCenterHandler struct {
	ServerCenterDefaultHandler
}

func (this *ServerCenterHandler) GetServerName(ctx context.Context) string {
	return model.DefaultServerName
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
	list = serverCenterClient.PingCheckAddress(ctx, list)
	logrus.WithContext(ctx).WithFields(logrus.Fields{"list": list}).Info("加载server_center地址")
	addresses = list
	return nil
}
func (this *ServerCenterHandler) GetDefaultConf(ctx context.Context) string {
	var config model.Config
	return util.ToYamlString(config)
}
