package main

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/config"
	"github.com/cellargalaxy/server_center/controller"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(config.Config.LogLevel)
	util.InitLog(util.GetServerNameWithPanic())
}

/**
server_name=server_center
log_level=4
mysql_dsn=root:123456@tcp(127.0.0.1:3306)/server_center?charset=utf8mb4&parseTime=True&loc=Local&tls=preferred
show_sql=false
secret=secret
*/
func main() {
	err := controller.Controller()
	if err != nil {
		panic(err)
	}
}
