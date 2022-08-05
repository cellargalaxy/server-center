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

server_name=server_center;server_center_address=http://127.0.0.1:7557;server_center_secret=secret
*/

type ServerCenterHandler struct {
	sdk.ServerCenterDefaultHandler
}

func (this *ServerCenterHandler) ListAddress(ctx context.Context) []string {
	return []string{"http://127.0.0.1:7557"}
}
func (this *ServerCenterHandler) GetSecret(ctx context.Context) string {
	return "secret"
}
func (this *ServerCenterHandler) GetServerName(ctx context.Context) string {
	return model.DefaultServerName
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
	client.StartWithInitConf(ctx)
	time.Sleep(time.Hour)
}

func TestAddEvent(test *testing.T) {
	for i := 0; i < 1000000000; i++ {
		ctx := util.GenCtx()
		sdk.AddEvent(ctx, "group", "name", 123, "data")
	}
	time.Sleep(time.Second * 3)
}

func TestPing(test *testing.T) {
	ctx := util.GenCtx()
	client, err := sdk.NewDefaultServerCenterClient(ctx, &ServerCenterHandler{})
	if err != nil {
		test.Error(err)
		test.FailNow()
	}
	object, err := client.Ping(ctx, "http://127.0.0.1:7557")
	test.Logf("object: %+v\r\n", util.ToJsonIndentString(object))
	if err != nil {
		test.Error(err)
		test.FailNow()
	}
}

func TestGetRemoteLastServerConf(test *testing.T) {
	ctx := util.GenCtx()
	client, err := sdk.NewDefaultServerCenterClient(ctx, &ServerCenterHandler{})
	if err != nil {
		test.Error(err)
		test.FailNow()
	}
	object, err := client.GetRemoteLastServerConf(ctx)
	test.Logf("object: %+v\r\n", util.ToJsonIndentString(object))
	if err != nil {
		test.Error(err)
		test.FailNow()
	}
}
