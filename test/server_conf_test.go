package test

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/cellargalaxy/server_center/service/db"
	"testing"
)

func TestAddServerConf(test *testing.T) {
	ctx := util.GenCtx()
	var object model.ServerConf
	object.ServerName = "server_center"
	object.Version = 10
	object.Remark = "remark"
	object.ConfText = `log_level: info
mysql_dsn: root:123456@tcp(127.0.0.1:3306)/server_center?charset=utf8mb4&parseTime=True&loc=Local&tls=preferred
show_sql: false
secret: secret_secret`
	response, err := db.AddServerConf(ctx, object)
	test.Logf("response: %+v\r\n", util.ToJsonIndentString(response))
	if err != nil {
		test.Error(err)
		test.FailNow()
	}
}

func TestListServerConf(test *testing.T) {
	ctx := util.GenCtx()
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
