package controller

import (
	"context"
	"github.com/cellargalaxy/server_center/model"
	"github.com/cellargalaxy/server_center/service/service"
)

func AddServerConf(ctx context.Context, request model.AddServerConfRequest) (*model.AddServerConfResponse, error) {
	object, err := service.AddServerConf(ctx, request.ServerConf)
	return &model.AddServerConfResponse{Conf: object}, err
}

func RemoveServerConf(ctx context.Context, request model.RemoveServerConfRequest) (model.RemoveServerConfResponse, error) {
	object, err := service.RemoveServerConf(ctx, request.ServerConfInquiry)
	return model.RemoveServerConfResponse{Conf: object}, err
}

func GetLastServerConf(ctx context.Context, request model.GetLastServerConfRequest) (*model.GetLastServerConfResponse, error) {
	object, err := service.GetLastServerConf(ctx, request.ServerConfInquiry)
	if object == nil {
		return nil, nil
	}
	return &model.GetLastServerConfResponse{Conf: object}, err
}

func ListServerConf(ctx context.Context, request model.ListServerConfRequest) (*model.ListServerConfResponse, error) {
	object, err := service.ListServerConf(ctx, request.ServerConfInquiry)
	return &model.ListServerConfResponse{List: object}, err
}

func ListAllServerName(ctx context.Context, request model.ListAllServerNameRequest) (*model.ListAllServerNameResponse, error) {
	object, err := service.ListAllServerName(ctx)
	return &model.ListAllServerNameResponse{List: object}, err
}
