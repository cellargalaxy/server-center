package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/cellargalaxy/server_center/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InsertServerConf(ctx context.Context, object model.ServerConfModel) (model.ServerConfModel, error) {
	err := getDb(ctx).Create(&object).Error
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Warn("插入服务配置，异常")
		return object, fmt.Errorf("插入服务配置，异常: %+v", err)
	}
	return object, nil
}

func serverConfWhere(where *gorm.DB, inquiry model.ServerConfInquiry) *gorm.DB {
	if inquiry.Id > 0 {
		where = where.Where("id = ?", inquiry.Id)
	}
	if inquiry.ServerName != "" {
		where = where.Where("server_name = ?", inquiry.ServerName)
	}
	if inquiry.Version > 0 {
		where = where.Where("version = ?", inquiry.Version)
	}
	if inquiry.CurrentVersion > 0 {
		where = where.Where("version > ?", inquiry.CurrentVersion)
	}
	return where
}

func SelectLastServerConf(ctx context.Context, inquiry model.ServerConfInquiry) (*model.ServerConfModel, error) {
	where := getDb(ctx)
	where = serverConfWhere(where, inquiry)
	where = where.Order("version desc")

	var object model.ServerConfModel
	err := where.Take(&object).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"inquiry": inquiry}).Warn("查询最新服务配置，不存在")
		return nil, nil
	}
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"inquiry": inquiry, "err": err}).Error("查询最新服务配置，异常")
		return nil, fmt.Errorf("查询最新服务配置，异常: %+v", err)
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{"object": object}).Debug("查询最新服务配置，完成")
	return &object, err
}

func SelectSomeServerConf(ctx context.Context, inquiry model.ServerConfInquiry) ([]model.ServerConfModel, error) {
	where := getDb(ctx)
	where = serverConfWhere(where, inquiry)

	var list []model.ServerConfModel
	err := where.Find(&list).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"inquiry": inquiry}).Warn("查询服务配置列表，不存在")
		return nil, nil
	}
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"inquiry": inquiry, "err": err}).Error("查询服务配置列表，异常")
		return nil, fmt.Errorf("查询服务配置列表，异常: %+v", err)
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{"len(list)": len(list)}).Debug("查询服务配置列表，完成")
	return list, err
}

func SelectServerConfDistinctServerName(ctx context.Context, inquiry model.ServerConfInquiry) ([]model.ServerConfModel, error) {
	where := getDb(ctx)
	where = where.Select("distinct server_name")
	where = serverConfWhere(where, inquiry)

	var list []model.ServerConfModel
	err := where.Find(&list).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"inquiry": inquiry}).Warn("查询服务配置列表，不存在")
		return nil, nil
	}
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"inquiry": inquiry, "err": err}).Error("查询服务配置列表，异常")
		return nil, fmt.Errorf("查询服务配置列表，异常: %+v", err)
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{"len(list)": len(list)}).Debug("查询服务配置列表，完成")
	return list, err
}
