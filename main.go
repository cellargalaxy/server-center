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
export server_name=server_center
export log_level=4
export mysql_dsn=root:123456@tcp(127.0.0.1:3306)/server_center?charset=utf8mb4&parseTime=True&loc=Local&tls=preferred
export show_sql=false
export secret=secret
*/
func main() {
	err := controller.Controller()
	if err != nil {
		panic(err)
	}
}
