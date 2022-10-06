package controller

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/cellargalaxy/server_center/service/controller"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func addEvent(ctx *gin.Context) {
	var request model.AddEventRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("插入批量事件，请求参数解析异常")
		ctx.JSON(http.StatusOK, util.NewHttpResponseByErr(err))
		return
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{}).Info("插入批量事件")
	ctx.JSON(http.StatusOK, util.NewHttpResponse(controller.AddEvent(ctx, request)))
}

func removeEvent(ctx *gin.Context) {
	var request model.RemoveEventRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("删除事件，请求参数解析异常")
		ctx.JSON(http.StatusOK, util.NewHttpResponseByErr(err))
		return
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{}).Info("删除事件")
	ctx.JSON(http.StatusOK, util.NewHttpResponse(controller.RemoveEvent(ctx, request)))
}
