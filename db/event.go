package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InsertEvents(ctx context.Context, object []model.EventModel) ([]model.EventModel, error) {
	ctx = util.CopyCtx(ctx)
	if len(object) == 0 {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Warn("插入批量事件，列表为空")
		return object, nil
	}
	err := getDb(ctx).Create(&object).Error
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("插入批量事件，异常")
		return object, fmt.Errorf("插入批量事件，异常: %+v", err)
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{}).Info("插入批量事件，完成")
	return object, nil
}

func DeleteEvent(ctx context.Context, inquiry model.EventInquiry) error {
	var where *gorm.DB
	where = whereEvent(ctx, where, inquiry)
	if where == nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"inquiry": inquiry}).Warn("删除事件，删除条件为空")
		return fmt.Errorf("删除事件，删除条件为空")
	}

	err := where.Delete(&inquiry.EventModel).Error
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Warn("删除事件，异常")
		return fmt.Errorf("删除事件，异常: %+v", err)
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{}).Info("删除事件，完成")
	return nil
}

func whereEvent(ctx context.Context, where *gorm.DB, inquiry model.EventInquiry) *gorm.DB {
	if inquiry.LogId > 0 {
		where = getWhere(ctx, where).Where("log_id = ?", inquiry.LogId)
	}
	if inquiry.ServerName != "" {
		where = getWhere(ctx, where).Where("server_name = ?", inquiry.ServerName)
	}
	if inquiry.Ip != "" {
		where = getWhere(ctx, where).Where("ip = ?", inquiry.Ip)
	}
	if inquiry.Group != "" {
		where = getWhere(ctx, where).Where("event_group = ?", inquiry.Group)
	}
	if inquiry.Name != "" {
		where = getWhere(ctx, where).Where("event_name = ?", inquiry.Name)
	}
	if inquiry.EndCreateTime.Unix() > 0 {
		where = getWhere(ctx, where).Where("create_time < ?", inquiry.EndCreateTime)
	}
	return where
}

func SelectSomeEvent(ctx context.Context, inquiry model.EventInquiry) ([]model.EventModel, error) {
	where := getDb(ctx)
	where = whereEvent(ctx, where, inquiry)
	where = where.Order("create_time desc")
	if inquiry.Offset > 0 {
		where = where.Offset(inquiry.Offset)
	}
	if inquiry.Limit > 0 {
		where = where.Offset(inquiry.Limit)
	}

	var object []model.EventModel
	err := where.Find(&object).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Warn("查询事件，不存在")
		return nil, nil
	}
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("查询事件，异常")
		return nil, fmt.Errorf("查询事件，异常: %+v", err)
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{"len(object)": len(object)}).Debug("查询事件，完成")
	return object, err
}
