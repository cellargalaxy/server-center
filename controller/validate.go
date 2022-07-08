package controller

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/config"
	"github.com/gin-gonic/gin"
)

func claims(ctx *gin.Context) {
	util.ClaimsHttp(ctx, config.Config.Secret)
}
func validate(ctx *gin.Context) {
	util.ValidateHttp(ctx, config.Config.Secret)
}
