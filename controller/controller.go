package controller

import (
	"fmt"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/cellargalaxy/server_center/static"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Controller() error {
	engine := gin.Default()
	engine.Use(util.GinLogId)
	engine.Use(util.GinLog)
	engine.GET("/ping", util.Ping)
	engine.POST("/ping", validate, util.Ping)

	engine.Use(staticCache)
	engine.StaticFS("/static", http.FS(static.StaticFile))

	engine.POST(model.AddServerConfPath, validate, addServerConf)
	engine.POST(model.RemoveServerConfPath, validate, removeServerConf)
	engine.GET(model.GetLastServerConfPath, validate, getLastServerConf)
	engine.GET(model.ListServerConfPath, validate, listServerConf)
	engine.GET(model.ListAllServerNamePath, validate, listAllServerName)

	err := engine.Run(model.ListenAddress)
	if err != nil {
		panic(fmt.Errorf("web服务启动，异常: %+v", err))
	}
	return nil
}

func staticCache(c *gin.Context) {
	if strings.HasPrefix(c.Request.RequestURI, "/static") {
		c.Header("Cache-Control", "max-age=86400")
	}
}
