package db

import (
	"context"
	common_model "github.com/cellargalaxy/go_common/model"
	"github.com/cellargalaxy/server_center/db"
	"github.com/cellargalaxy/server_center/model"
	"github.com/sirupsen/logrus"
)

func AddManyEvent(ctx context.Context, list []model.Event) {
	left := 0
	right := common_model.DbMaxBatchAddLength
	for true {
		if left >= len(list) {
			break
		}
		if right > len(list) {
			right = len(list)
		}
		addEvents(ctx, list[left:right])
		left = right
		right += common_model.DbMaxBatchAddLength
	}
}

func addEvents(ctx context.Context, object []model.Event) error {
	list := make([]model.EventModel, 0, len(object))
	for i := range object {
		if object[i].Group == "" {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"object": object[i]}).Error("插入批量事件，EventGroup非法")
			continue
		}
		if object[i].Name == "" {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"object": object[i]}).Error("插入批量事件，EventName非法")
			continue
		}
		list = append(list, model.EventModel{Event: object[i]})
	}
	_, err := db.InsertEvents(ctx, list)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"list": list}).Error("插入批量事件，异常")
	}
	return err
}

func ListEvent(ctx context.Context, inquiry model.EventInquiry) ([]model.EventModel, error) {
	return db.SelectSomeEvent(ctx, inquiry)
}

func RemoveEvent(ctx context.Context, inquiry model.EventInquiry) error {
	return db.DeleteEvent(ctx, inquiry)
}
