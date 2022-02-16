package main

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/config"
	"github.com/cellargalaxy/server_center/controller"
	"github.com/sirupsen/logrus"
)

func init() {
	ctx := util.CreateLogCtx()
	logrus.SetLevel(config.Config.LogLevel)
	util.InitDefaultLog()
	config.Init(ctx)
}

/**
export server_name=server_center
export mysql_dsn=root:123456@tcp(172.17.0.2:3306)/server_center?charset=utf8mb4&parseTime=True&loc=Local&tls=preferred

server_name=server_center;mysql_dsn=root:123456@tcp(127.0.0.1:3306)/server_center?charset=utf8mb4&parseTime=True&loc=Local&tls=preferred
*/
func main() {
	err := controller.Controller()
	if err != nil {
		panic(err)
	}
}
