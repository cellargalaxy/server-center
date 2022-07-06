package test

import (
	"context"
	"fmt"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/cellargalaxy/server_center/sdk"
	"testing"
	"time"
)

/**
export server_name=server_center
export server_center_address=http://127.0.0.1:7557
export server_center_secret=secret_secret

server_name=server_center;server_center_address=http://127.0.0.1:7557;server_center_secret=secret_secret
*/

type ServerCenterHandler struct {
}

func (this *ServerCenterHandler) GetAddress(ctx context.Context) string {
	return sdk.GetEnvServerCenterAddress(ctx)
}
func (this *ServerCenterHandler) GetSecret(ctx context.Context) string {
	return sdk.GetEnvServerCenterSecret(ctx)
}
func (this *ServerCenterHandler) GetServerName(ctx context.Context) string {
	return sdk.GetEnvServerName(ctx)
}
func (this *ServerCenterHandler) GetInterval(ctx context.Context) time.Duration {
	return 5 * time.Second
}
func (this *ServerCenterHandler) ParseConf(ctx context.Context, object model.ServerConfModel) error {
	fmt.Printf("解析配置: \n%+v\n", object.ConfText)
	return nil
}
func (this *ServerCenterHandler) GetDefaultConf(ctx context.Context) string {
	return ""
}

func TestGetAndParseLastServerConf(test *testing.T) {
	ctx := util.GenCtx()
	client, err := sdk.NewDefaultServerCenterClient(ctx, &ServerCenterHandler{})
	if err != nil {
		test.Error(err)
		test.FailNow()
	}
	response, err := client.GetAndParseLastServerConf(ctx)
	test.Logf("response: %+v\r\n", util.ToJsonIndentString(response))
	if err != nil {
		test.Error(err)
		test.FailNow()
	}
}

func TestStartConfWithInitConf(test *testing.T) {
	ctx := util.GenCtx()
	client, err := sdk.NewDefaultServerCenterClient(ctx, &ServerCenterHandler{})
	if err != nil {
		test.Error(err)
		test.FailNow()
	}
	response, err := client.StartConfWithInitConf(ctx)
	test.Logf("response: %+v\r\n", util.ToJsonIndentString(response))
	if err != nil {
		test.Error(err)
		test.FailNow()
	}
	time.Sleep(time.Hour)
}
