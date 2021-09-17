package test

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/cellargalaxy/server_center/service/db"
	"testing"
)

func TestAddServerConf(test *testing.T) {
	ctx := util.CreateLogCtx()
	var object model.ServerConf
	object.ServerName = "server_center"
	object.Version = 4
	object.Remark = "remark"
	object.ConfText = `;Debug:5 Info:4 Warn:3 Error:2
log_level = 4
retry = 3
timeout = "3s"
sleep = "3s"
secret = "secret"
wx_app_id = ""
wx_app_secret = ""
wx_common_temp_id = ""
wx_common_tag_id = 0
tg_token = ""
tg_chat_id = 0`
	response, err := db.AddServerConf(ctx, object)
	test.Logf("response: %+v\r\n", util.ToJsonIndentString(response))
	if err != nil {
		test.Error(err)
		test.FailNow()
	}
}

func TestListServerConf(test *testing.T) {
	ctx := util.CreateLogCtx()
	var object model.ServerConfInquiry
	object.ServerName = "server_center"
	object.Version = 2
	response, err := db.ListServerConf(ctx, object)
	test.Logf("response: %+v\r\n", util.ToJsonIndentString(response))
	if err != nil {
		test.Error(err)
		test.FailNow()
	}
}

func TestListAllServerName(test *testing.T) {
	ctx := util.CreateLogCtx()
	response, err := db.ListAllServerName(ctx)
	test.Logf("response: %+v\r\n", util.ToJsonIndentString(response))
	if err != nil {
		test.Error(err)
		test.FailNow()
	}
}
