package db

import (
	"context"
	"fmt"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/mmm/db"
	"github.com/cellargalaxy/mmm/service/msg"
	"github.com/cellargalaxy/server-center/model"
	"github.com/sirupsen/logrus"
)

func AddTrackEvent(ctx context.Context, object model.ServerConf) (model.ServerConfModel, error) {
	var event model.TrackEventModel
	event.TrackEvent = object
	event.LogId = util.GetLogId(ctx)
	if event.Key == "" {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"object": event}).Error("插入埋点，Key为空")
		return event, fmt.Errorf("插入埋点，Key为空")
	}
	event.Data = util.ToJsonString(event.DataMap)
	event, err := db.InsertTrackEvent(ctx, event)
	msg.SendDbErr(ctx, "AddTrackEvent", err)
	return event, err
}

func ListTrackEvent(ctx context.Context, inquiry model.ServerConfInquiry) ([]model.ServerConfModel, error) {
	list, err := db.SelectSomeTrackEvent(ctx, inquiry)
	for i := range list {
		if list[i].DataMap == nil {
			list[i].DataMap = make(map[string]interface{})
		}
		util.UnmarshalJsonString(list[i].Data, &list[i].DataMap)
	}
	return list, err
}