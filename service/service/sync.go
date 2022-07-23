package service

import (
	"context"
	"fmt"
	"github.com/cellargalaxy/server_center/config"
	"github.com/cellargalaxy/server_center/model"
	"github.com/cellargalaxy/server_center/sdk"
	"github.com/sirupsen/logrus"
	"time"
)

type PullSyncHandler struct {
}

func (this *PullSyncHandler) ListAddress(ctx context.Context) []string {
	pullSyncHost := config.Config.PullSyncHost
	if pullSyncHost != "" {
		return []string{pullSyncHost}
	}
	return sdk.ListAddress(ctx)
}
func (this *PullSyncHandler) GetSecret(ctx context.Context) string {
	secret := config.Config.PullSyncSecret
	if secret != "" {
		return secret
	}
	return sdk.GetSecret(ctx)
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

func PullSync(ctx context.Context) error {
	var handler PullSyncHandler
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
		AddServerConf(ctx, conf.ServerConf)
	}
	return nil
}
