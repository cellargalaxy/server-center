package service

import (
	"context"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/config"
	"github.com/cellargalaxy/server_center/model"
	"github.com/cellargalaxy/server_center/service/db"
	"github.com/sirupsen/logrus"
)

func initEvent(ctx context.Context) {
	flushEventSync()
}

var eventChan = make(chan []model.Event, 1000)

func AddEventsAsync(ctx context.Context, object []model.Event) {
	go func() {
		defer util.Defer(ctx, func(ctx context.Context, err interface{}, stack string) {
			if err != nil {
				logrus.WithContext(ctx).WithFields(logrus.Fields{"object": object, "err": err, "stack": stack}).Error("插入批量事件，异常")
			}
		})

		AddEvents(ctx, object)
	}()
}

func AddEvents(ctx context.Context, object []model.Event) {
	if len(object) == 0 {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Warn("插入批量事件，为空")
		return
	}
	claims := util.GetClaims(ctx)
	if claims != nil {
		for i := range object {
			if object[i].Ip == "" {
				object[i].Ip = claims.Ip
			}
			if object[i].ServerName == "" {
				object[i].ServerName = claims.ServerName
			}
			if object[i].LogId <= 0 {
				object[i].LogId = claims.LogId
			}
		}
	}
	eventChan <- object
}

func flushEventSync() {
	ctx := util.GenCtx()
	go func() {
		defer util.Defer(ctx, func(ctx context.Context, err interface{}, stack string) {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err, "stack": stack}).Error("刷入事件，退出")
			flushEventSync()
		})

		flushEvent(ctx)
	}()
}

func flushEvent(ctx context.Context) {
	for events := range eventChan {
		ctx = util.GenCtx()
		db.AddManyEvent(ctx, events)
	}
}

func getOldEvent(ctx context.Context, maxSave int) (*model.EventModel, error) {
	var inquiry model.EventInquiry
	inquiry.Offset = maxSave
	inquiry.Limit = 1
	object, err := db.ListEvent(ctx, inquiry)
	if err != nil {
		return nil, err
	}
	if len(object) == 0 {
		return nil, nil
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{"object[0]": object[0]}).Info("查询旧事件")
	return &object[0], nil
}

func ClearEvent(ctx context.Context) error {
	maxSave := config.Config.EventMaxSave
	if maxSave <= 0 {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Info("清理旧事件，不进行清理")
		return nil
	}
	object, err := getOldEvent(ctx, maxSave)
	if err != nil {
		return err
	}
	if object == nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Info("清理旧事件，无旧事件")
		return nil
	}
	var inquiry model.EventInquiry
	inquiry.EndCreatedAt = object.CreatedAt
	err = db.RemoveEvent(ctx, inquiry)
	return err
}
