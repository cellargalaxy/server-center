package controller

import (
	"context"
	"github.com/cellargalaxy/server_center/model"
	"github.com/cellargalaxy/server_center/service/db"
)

func AddServerConf(ctx context.Context, request model.AddServerConfRequest) (*model.AddServerConfResponse, error) {
	object, err := db.AddServerConf(ctx, request.ServerConf)
	return &model.AddServerConfResponse{Conf: &object}, err
}

func GetLastServerConf(ctx context.Context, request model.GetLastServerConfRequest) (*model.GetLastServerConfResponse, error) {
	object, err := db.GetLastServerConf(ctx, request.ServerConfInquiry)
	if object == nil {
		return nil, nil
	}
	return &model.GetLastServerConfResponse{Conf: object}, err
}

func ListServerConf(ctx context.Context, request model.ListServerConfRequest) (*model.ListServerConfResponse, error) {
	object, err := db.ListServerConf(ctx, request.ServerConfInquiry)
	return &model.ListServerConfResponse{List: object}, err
}

func ListAllServerName(ctx context.Context, request model.ListAllServerNameRequest) (*model.ListAllServerNameResponse, error) {
	object, err := db.ListAllServerName(ctx)
	names := make([]string, 0, len(object))
	for i := range object {
		names = append(names, object[i].ServerName)
	}
	return &model.ListAllServerNameResponse{List: names}, err
}
