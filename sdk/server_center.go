package sdk

import (
	"context"
	"fmt"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/sirupsen/logrus"
	"time"
)

var addresses []string
var secret string
var client *ServerCenterClient
var eventChan = make(chan model.Event, util.DbMaxBatchAddLength)

func initServerCenter(ctx context.Context) {
	var err error

	address := GetEnvServerCenterAddress(ctx)
	if address != "" {
		addresses = append(addresses, address)
	}
	secret = GetEnvServerCenterSecret(ctx)

	var handler ServerCenterHandler
	client, err = NewDefaultServerCenterClient(ctx, &handler)
	if err != nil {
		panic(err)
	}
	if client == nil {
		panic("创建serverCenterClient为空")
	}
	client.StartWithInitConf(ctx)

	flushEventAsync(ctx)
}

func ListAddress(ctx context.Context) []string {
	return addresses
}
func GetSecret(ctx context.Context) string {
	return secret
}

func AddEvent(ctx context.Context, group, name string, value float64, data interface{}) {
	var event model.Event
	event.LogId = util.GetLogId(ctx)
	event.ServerName = GetEnvServerName(ctx, "")
	event.Ip = util.GetIp()
	event.Group = group
	event.Name = name
	event.Value = value
	event.Data = fmt.Sprint(data)
	AddEventAsync(ctx, event)
}
func AddEventAsync(ctx context.Context, event model.Event) {
	go func() {
		defer util.Defer(ctx, func(ctx context.Context, err interface{}, stack string) {
			if err != nil {
				logrus.WithContext(ctx).WithFields(logrus.Fields{"event": event, "err": err, "stack": stack}).Error("插入事件，异常")
			}
		})
		addEvent(ctx, event)
	}()
}
func addEvent(ctx context.Context, event model.Event) {
	eventChan <- event
}
func flushEventAsync(ctx context.Context) {
	go func() {
		defer util.Defer(ctx, func(ctx context.Context, err interface{}, stack string) {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err, "stack": stack}).Error("插入事件，退出")
			flushEventAsync(ctx)
		})

		flushEvent(ctx)
	}()
}
func flushEvent(ctx context.Context) {
	list := make([]model.Event, 0, util.DbMaxBatchAddLength)

	defer util.Defer(ctx, func(ctx context.Context, err interface{}, stack string) {
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err, "stack": stack}).Error("插入事件，异常")
		}
		if len(list) == 0 || client == nil {
			return
		}
		client.AddEvent(ctx, list)
	})

	for {
		ctx = util.ResetLogId(ctx)
		select {
		case event := <-eventChan:
			list = append(list, event)
			if len(list) < util.DbMaxBatchAddLength {
				continue
			}
			if client == nil {
				logrus.WithContext(ctx).WithFields(logrus.Fields{"list": list}).Error("插入事件，serverCenterClient为空")
			} else {
				client.AddEvent(ctx, list)
			}
			list = make([]model.Event, 0, util.DbMaxBatchAddLength)
		case <-time.After(time.Second):
			if len(list) == 0 {
				continue
			}
			if client == nil {
				logrus.WithContext(ctx).WithFields(logrus.Fields{"list": list}).Error("插入事件，serverCenterClient为空")
			} else {
				client.AddEvent(ctx, list)
			}
			list = make([]model.Event, 0, util.DbMaxBatchAddLength)
		case <-ctx.Done():
			return
		}
	}
}

type ServerCenterDefaultHandler struct {
	intervalIndex int
}

func (this *ServerCenterDefaultHandler) ListAddress(ctx context.Context) []string {
	return ListAddress(ctx)
}
func (this *ServerCenterDefaultHandler) GetSecret(ctx context.Context) string {
	return GetSecret(ctx)
}
func (this *ServerCenterDefaultHandler) GetInterval(ctx context.Context) time.Duration {
	intervals := []time.Duration{time.Second * 2, time.Second * 4, time.Second * 8, time.Second * 16, time.Second * 32, time.Minute * 10}
	index := this.intervalIndex % len(intervals)
	this.intervalIndex = this.intervalIndex + 1
	return intervals[index]
}

type ServerCenterHandler struct {
	ServerCenterDefaultHandler
}

func (this *ServerCenterHandler) GetServerName(ctx context.Context) string {
	return model.DefaultServerName
}
func (this *ServerCenterHandler) ParseConf(ctx context.Context, object model.ServerConfModel) error {
	var config model.Config
	err := util.UnmarshalYamlString(object.ConfText, &config)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("反序列化server_center配置异常")
		return err
	}
	secret = config.Secret

	list := addresses
	address := GetEnvServerCenterAddress(ctx)
	if address != "" {
		list = append(list, address)
	}
	list = append(list, config.Addresses...)
	list = util.DistinctString(ctx, list)
	list = client.PingCheckAddress(ctx, list)
	logrus.WithContext(ctx).WithFields(logrus.Fields{"list": list}).Info("加载server_center地址")
	addresses = list
	return nil
}
func (this *ServerCenterHandler) GetDefaultConf(ctx context.Context) string {
	var config model.Config
	return util.ToYamlString(config)
}
