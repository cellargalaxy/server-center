package service

import (
	"context"
	"github.com/cellargalaxy/server_center/config"
	"github.com/cellargalaxy/server_center/model"
	"github.com/cellargalaxy/server_center/service/db"
	"github.com/sirupsen/logrus"
)

func AddServerConf(ctx context.Context, object model.ServerConf) (model.ServerConfModel, error) {
	return db.AddServerConf(ctx, object)
}

func RemoveServerConf(ctx context.Context, inquiry model.ServerConfInquiry) (model.ServerConfModel, error) {
	return db.RemoveServerConf(ctx, inquiry)
}

func GetLastServerConf(ctx context.Context, inquiry model.ServerConfInquiry) (*model.ServerConfModel, error) {
	object, err := db.GetLastServerConf(ctx, inquiry)
	return object, err
}

func ListServerConf(ctx context.Context, inquiry model.ServerConfInquiry) ([]model.ServerConfModel, error) {
	list, err := db.ListServerConf(ctx, inquiry)
	return list, err
}

func ListAllServerName(ctx context.Context) ([]string, error) {
	var inquiry model.ServerConfInquiry
	list, err := db.ListServerName(ctx, inquiry)
	return list, err
}

func getOldServerConfByServerName(ctx context.Context, serverName string, maxSave int) (*model.ServerConfModel, error) {
	var inquiry model.ServerConfInquiry
	inquiry.ServerName = serverName
	inquiry.Offset = maxSave
	inquiry.Limit = 1
	object, err := ListServerConf(ctx, inquiry)
	if err != nil {
		return nil, err
	}
	if len(object) == 0 {
		return nil, nil
	}
	return &object[0], nil
}

func listAllOldServerConf(ctx context.Context, maxSave int) ([]model.ServerConfModel, error) {
	names, err := ListAllServerName(ctx)
	if err != nil {
		return nil, err
	}
	confs := make([]model.ServerConfModel, 0, len(names))
	for i := range names {
		conf, err := getOldServerConfByServerName(ctx, names[i], maxSave)
		if conf == nil || err != nil {
			continue
		}
		confs = append(confs, *conf)
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{"confs": confs}).Info("查询旧配置")
	return confs, nil
}

func ClearConfig(ctx context.Context) error {
	maxSave := config.Config.ClearConfigSave
	if maxSave <= 0 {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Info("清理旧配置，不进行清理")
		return nil
	}
	confs, err := listAllOldServerConf(ctx, maxSave)
	if err != nil {
		return err
	}
	if len(confs) == 0 {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Info("清理旧配置，无旧配置")
		return nil
	}
	for i := range confs {
		var inquiry model.ServerConfInquiry
		inquiry.ServerName = confs[i].ServerName
		inquiry.EndCreatedAt = confs[i].CreatedAt
		_, removeErr := RemoveServerConf(ctx, inquiry)
		if removeErr != nil {
			err = removeErr
		}
	}
	return err
}
