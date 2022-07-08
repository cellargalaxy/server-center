package controller

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/cellargalaxy/server_center/service/controller"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func addServerConf(ctx *gin.Context) {
	var request model.AddServerConfRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("插入服务配置，请求参数解析异常")
		ctx.JSON(http.StatusOK, util.CreateResponseByErr(err))
		return
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{}).Info("插入服务配置")
	ctx.JSON(http.StatusOK, util.CreateResponse(controller.AddServerConf(ctx, request)))
}

func removeServerConf(ctx *gin.Context) {
	var request model.RemoveServerConfRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"request": request, "err": err}).Error("删除服务配置，请求参数解析异常")
		ctx.JSON(http.StatusOK, util.CreateResponseByErr(err))
		return
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{"request": request}).Info("删除服务配置")
	ctx.JSON(http.StatusOK, util.CreateResponse(controller.RemoveServerConf(ctx, request)))
}

func getLastServerConf(ctx *gin.Context) {
	var request model.GetLastServerConfRequest
	err := ctx.BindQuery(&request)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"request": request, "err": err}).Error("查询最新服务配置，请求参数解析异常")
		ctx.JSON(http.StatusOK, util.CreateResponseByErr(err))
		return
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{"request": request}).Info("查询最新服务配置")
	ctx.JSON(http.StatusOK, util.CreateResponse(controller.GetLastServerConf(ctx, request)))
}

func listServerConf(ctx *gin.Context) {
	var request model.ListServerConfRequest
	err := ctx.BindQuery(&request)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"request": request, "err": err}).Error("查询服务配置列表，请求参数解析异常")
		ctx.JSON(http.StatusOK, util.CreateResponseByErr(err))
		return
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{"request": request}).Info("查询服务配置列表")
	ctx.JSON(http.StatusOK, util.CreateResponse(controller.ListServerConf(ctx, request)))
}

func listAllServerName(ctx *gin.Context) {
	var request model.ListAllServerNameRequest
	err := ctx.BindQuery(&request)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"request": request, "err": err}).Error("查询服务配置列表，请求参数解析异常")
		ctx.JSON(http.StatusOK, util.CreateResponseByErr(err))
		return
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{"request": request}).Info("查询服务配置列表")
	ctx.JSON(http.StatusOK, util.CreateResponse(controller.ListAllServerName(ctx, request)))
}
