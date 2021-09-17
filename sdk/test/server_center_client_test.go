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
export server_center_secret=secret
*/

type ServerCenterHandler struct {
}

func (this *ServerCenterHandler) ParseConf(ctx context.Context, object model.ServerConfModel) error {
	fmt.Printf("解析配置: \n%+v\n", object.ConfText)
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
