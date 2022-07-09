package service

import (
	"context"
	"fmt"
	"github.com/cellargalaxy/server_center/model"
	"github.com/cellargalaxy/server_center/sdk"
	"github.com/cellargalaxy/server_center/service/db"
	"github.com/sirupsen/logrus"
	"time"
)

type PullSyncHandler struct {
	address string
	secret  string
}

func (this *PullSyncHandler) ListAddress(ctx context.Context) []string {
	return []string{this.address}
}
func (this *PullSyncHandler) GetSecret(ctx context.Context) string {
	return this.secret
}
func (this *PullSyncHandler) GetServerName(ctx context.Context) string {
	return model.DefaultServerName
}
func (this *PullSyncHandler) GetInterval(ctx context.Context) time.Duration {
	return 5 * time.Minute
}
func (this *PullSyncHandler) ParseConf(ctx context.Context, object model.ServerConfModel) error {
	return nil
}
func (this *PullSyncHandler) GetDefaultConf(ctx context.Context) string {
	return ""
}

func PullSync(ctx context.Context, address, secret string) error {
	var handler PullSyncHandler
	handler.address = address
	handler.secret = secret
	client, err := sdk.NewDefaultServerCenterClient(ctx, &handler)
	if err != nil {
		return err
	}
	if client == nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("同步拉服务配置，ServerCenterClient为空")
		return fmt.Errorf("同步拉服务配置，ServerCenterClient为空")
	}
	names, err := client.ListAllServerName(ctx)
	if err != nil {
		return err
	}
	for i := range names {
		conf, err := client.GetRemoteLastServerConfByServerName(ctx, names[i])
		if conf == nil || err != nil {
			continue
		}
		db.AddServerConf(ctx, conf.ServerConf)
	}
	return nil
}
