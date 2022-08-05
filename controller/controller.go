package controller

import (
	"fmt"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/cellargalaxy/server_center/static"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Controller() error {
	engine := gin.Default()
	engine.Use(claims)
	engine.Use(util.GinLog)

	debug := engine.Group(util.DebugPath, validate)
	pprof.RouteRegister(debug, util.PprofPath)

	engine.GET(util.PingPath, util.Ping)
	engine.POST(util.PingPath, validate, util.Ping)

	engine.Use(staticCache)
	engine.StaticFS(util.StaticPath, http.FS(static.StaticFile))

	engine.POST(model.AddServerConfPath, validate, addServerConf)
	engine.POST(model.RemoveServerConfPath, validate, removeServerConf)
	engine.GET(model.GetLastServerConfPath, validate, getLastServerConf)
	engine.GET(model.ListServerConfPath, validate, listServerConf)
	engine.GET(model.ListAllServerNamePath, validate, listAllServerName)
	engine.POST(model.AddEventPath, validate, addEvent)

	err := engine.Run(model.ListenAddress)
	if err != nil {
		panic(fmt.Errorf("web服务启动，异常: %+v", err))
	}
	return nil
}

func staticCache(c *gin.Context) {
	if strings.HasPrefix(c.Request.RequestURI, util.StaticPath) {
		c.Header("Cache-Control", "max-age=86400")
	}
}
