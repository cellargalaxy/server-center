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

type ServerCenterHandler struct {
}

func (this *ServerCenterHandler) ParseConf(ctx context.Context, object model.ServerConfModel) error {
	fmt.Printf("解析配置: %+v\n", util.ToJsonIndentString(object))
	return nil
}

func TestGetAndParseLastServerConf(test *testing.T) {
	ctx := util.CreateLogCtx()
	util.InitLog("server_center.log")
	client, err := sdk.NewDefaultServerCenterClient(&ServerCenterHandler{})
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

func TestStart(test *testing.T) {
	ctx := util.CreateLogCtx()
	util.InitLog("server_center.log")
	client, err := sdk.NewDefaultServerCenterClient(&ServerCenterHandler{})
	if err != nil {
		test.Error(err)
		test.FailNow()
	}
	response, err := client.Start(ctx)
	test.Logf("response: %+v\r\n", util.ToJsonIndentString(response))
	if err != nil {
		test.Error(err)
		test.FailNow()
	}
	time.Sleep(time.Hour)
}
