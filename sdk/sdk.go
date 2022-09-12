package sdk

import (
	"context"
	"fmt"
	common_model "github.com/cellargalaxy/go_common/model"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

const (
	ServerCenterAddressEnvKey = "server_center_address"
	ServerCenterSecretEnvKey  = "server_center_secret"
)

func init() {
	ctx := util.GenCtx()
	initServerCenter(ctx)
}

func GetEnvServerCenterAddress(ctx context.Context) string {
	return util.GetEnvString(ServerCenterAddressEnvKey, "")
}
func GetEnvServerCenterSecret(ctx context.Context) string {
	return util.GetEnvString(ServerCenterSecretEnvKey, "")
}
func GetEnvServerName(ctx context.Context, defaultServerName string) string {
	serverName := util.GetServerName()
	if serverName == "" {
		serverName = defaultServerName
	}
	return serverName
}

type ServerCenterHandlerInter interface {
	ListAddress(ctx context.Context) []string
	GetSecret(ctx context.Context) string
	GetServerName(ctx context.Context) string
	GetInterval(ctx context.Context) time.Duration
	ParseConf(ctx context.Context, object model.ServerConfModel) error
	GetDefaultConf(ctx context.Context) string
}

func GenName(ctx context.Context, handler ServerCenterHandlerInter) string {
	if handler == nil {
		return "ServerCenterClient"
	}
	return fmt.Sprintf("ServerCenterClient_%s", handler.GetServerName(ctx))
}

func NewDefaultServerCenterClient(ctx context.Context, handler ServerCenterHandlerInter) (*ServerCenterClient, error) {
	var err error
	client, err = NewServerCenterClient(ctx, util.TimeoutDefault, util.TryDefault, util.GetHttpClient(), handler)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewServerCenterClient(ctx context.Context, timeout time.Duration, try int, httpClient *resty.Client, handler ServerCenterHandlerInter) (*ServerCenterClient, error) {
	if handler == nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Error("创建ServerCenterClient，handler为空")
		return nil, fmt.Errorf("创建ServerCenterClient，handler为空")
	}

	name := GenName(ctx, handler)
	client = &ServerCenterClient{timeout: timeout, try: try, httpClient: httpClient, handler: handler, name: name}
	return client, nil
}

type ServerCenterClient struct {
	timeout    time.Duration
	try        int
	httpClient *resty.Client
	handler    ServerCenterHandlerInter
	name       string

	lock sync.Mutex
	pool *util.SingleGoPool

	conf model.ServerConfModel
}

func (this *ServerCenterClient) StartWithInitConf(ctx context.Context) error {
	for {
		_, err := this.GetAndParseLastServerConf(ctx)
		if err == nil {
			break
		}
		util.SleepWare(ctx, time.Second)
		if util.CtxDone(ctx) {
			return nil
		}
	}
	return this.Start(ctx)
}
func (this *ServerCenterClient) Start(ctx context.Context) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.pool != nil && !this.pool.IsClosed(ctx) {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"name": this.getName(ctx)}).Warn("已启动")
		return nil
	}

	var err error
	this.pool, err = util.NewDaemonSingleGoPool(ctx, this.getName(ctx), time.Second, this.flushServerConf)
	if err != nil {
		util.ReleasePool(ctx, this.pool)
		this.pool = nil
		return err
	}

	return nil
}
func (this *ServerCenterClient) flushServerConf(ctx context.Context, pool *util.SingleGoPool) {
	defer util.Defer(func(err interface{}, stack string) {
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err, "stack": stack}).Warn("flushServerConf，异常")
		}
	})

	for {
		ctx := util.ResetLogId(ctx)
		this.GetAndParseLastServerConf(ctx)
		util.SleepWare(ctx, this.handler.GetInterval(ctx))
		if util.CtxDone(ctx) {
			return
		}
	}
}

var confLock sync.Mutex

func (this *ServerCenterClient) GetAndParseLastServerConf(ctx context.Context) (*model.ServerConfModel, error) {
	confLock.Lock()
	defer confLock.Unlock()

	object, err := this.getLastServerConf(ctx)
	if err != nil {
		return nil, err
	}
	if object == nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"server_name": this.handler.GetServerName(ctx), "current_version": this.conf.Version}).Error("查询并解析最新服务配置，服务配置为空")
		return nil, fmt.Errorf("查询并解析最新服务配置，服务配置为空")
	}
	err = this.handler.ParseConf(ctx, *object)
	if err != nil {
		return object, err
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{"server_name": this.handler.GetServerName(ctx), "current_version": this.conf.Version, "last_version": object.Version}).Info("查询并解析最新服务配置，查询服务配置")
	if object.Version <= this.conf.Version {
		return nil, nil
	}
	if object.ConfText == this.conf.ConfText {
		return nil, nil
	}
	this.conf = *object
	this.saveLocalFileServerConf(ctx, this.conf.ConfText)
	return object, nil
}
func (this *ServerCenterClient) getLastServerConf(ctx context.Context) (*model.ServerConfModel, error) {
	if this.getAddress(ctx) == "" {
		return this.getLocalFileServerConf(ctx)
	}
	object, err := this.GetRemoteLastServerConf(ctx)
	if err == nil {
		return object, nil
	}
	return this.getLocalFileServerConf(ctx)
}
func (this *ServerCenterClient) getLocalFileServerConf(ctx context.Context) (*model.ServerConfModel, error) {
	filePath, err := this.getLocalFilePath(ctx)
	if err != nil {
		return nil, err
	}
	if filePath == "" {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Warn("查询本地文件服务配置，filePath为空")
		return nil, nil
	}
	confText, err := util.ReadFileWithString(ctx, filePath, "")
	if err != nil {
		return nil, err
	}
	if confText == "" {
		confText = this.handler.GetDefaultConf(ctx)
		this.saveLocalFileServerConf(ctx, confText)
	}
	var serverConf model.ServerConfModel
	serverConf.ServerName = this.handler.GetServerName(ctx)
	serverConf.Version = this.conf.Version + 1 //本地文件更新了也能更新配置
	serverConf.ConfText = confText
	return &serverConf, nil
}
func (this *ServerCenterClient) saveLocalFileServerConf(ctx context.Context, confText string) error {
	filePath, err := this.getLocalFilePath(ctx)
	if err != nil {
		return err
	}
	if filePath == "" {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Warn("保存本地文件服务配置，filePath为空")
		return nil
	}
	return util.WriteFileWithString(ctx, filePath, confText)
}
func (this *ServerCenterClient) getLocalFilePath(ctx context.Context) (string, error) {
	serverName := this.handler.GetServerName(ctx)
	if serverName == "" {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Warn("查询本地文件服务配置，serverName为空")
		return "", nil
	}
	filePath := "resource/" + serverName + ".yml"
	logrus.WithContext(ctx).WithFields(logrus.Fields{"filePath": filePath}).Info("查询本地文件服务配置")
	return filePath, nil
}
func (this *ServerCenterClient) GetRemoteLastServerConf(ctx context.Context) (*model.ServerConfModel, error) {
	return this.GetRemoteLastServerConfByServerName(ctx, this.handler.GetServerName(ctx))
}
func (this *ServerCenterClient) GetRemoteLastServerConfByServerName(ctx context.Context, serverName string) (*model.ServerConfModel, error) {
	if serverName == "" {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Warn("查询最新服务配置，serverName为空")
		return nil, nil
	}
	ctx = util.RmReqId(ctx)
	type Response struct {
		common_model.HttpResponse
		Data model.GetLastServerConfResponse `json:"data"`
	}
	var response Response
	err := util.HttpApiWithTry(ctx, "查询最新服务配置", this.try, nil, &response, func() (*resty.Response, error) {
		response, err := this.httpClient.R().SetContext(ctx).
			SetHeader(this.genJWT(ctx)).
			SetQueryParam("server_name", serverName).
			Get(this.GetUrl(ctx, model.GetLastServerConfPath))
		return response, err
	})
	return response.Data.Conf, err
}

func (this *ServerCenterClient) ListAllServerName(ctx context.Context) ([]string, error) {
	if this.getAddress(ctx) == "" {
		return this.ListLocalAllServerName(ctx)
	}
	object, err := this.ListRemoteAllServerName(ctx)
	if len(object) > 0 && err == nil {
		return object, nil
	}
	return this.ListLocalAllServerName(ctx)
}
func (this *ServerCenterClient) ListLocalAllServerName(ctx context.Context) ([]string, error) {
	list := make([]string, 0, 1)
	serverName := this.handler.GetServerName(ctx)
	if serverName != "" {
		list = append(list, serverName)
	}
	return list, nil
}
func (this *ServerCenterClient) ListRemoteAllServerName(ctx context.Context) ([]string, error) {
	ctx = util.RmReqId(ctx)
	type Response struct {
		common_model.HttpResponse
		Data model.ListAllServerNameResponse `json:"data"`
	}
	var response Response
	err := util.HttpApiWithTry(ctx, "查询服务配置列表", this.try, nil, &response, func() (*resty.Response, error) {
		response, err := this.httpClient.R().SetContext(ctx).
			SetHeader(this.genJWT(ctx)).
			Get(this.GetUrl(ctx, model.ListAllServerNamePath))
		return response, err
	})
	return response.Data.List, err
}

func (this *ServerCenterClient) AddEvent(ctx context.Context, object []model.Event) error {
	if this.getAddress(ctx) == "" {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"object": object}).Warn("插入批量事件，地址为空")
		return nil
	}
	ctx = util.SetReqId(ctx)
	type Response struct {
		common_model.HttpResponse
		Data model.AddEventResponse `json:"data"`
	}
	var response Response
	err := util.HttpApiWithTry(ctx, "插入批量事件", this.try, nil, &response, func() (*resty.Response, error) {
		var req model.AddEventRequest
		req.List = object
		response, err := this.httpClient.R().SetContext(ctx).
			SetHeader(this.genJWT(ctx)).
			SetBody(req).
			Post(this.GetUrl(ctx, model.AddEventPath))
		return response, err
	})
	return err
}

func (this *ServerCenterClient) PingCheckAddress(ctx context.Context, addresses []string) []string {
	listChan := make(chan string, len(addresses))
	var wg sync.WaitGroup
	for i := range addresses {
		wg.Add(1)
		go func(address string) {
			defer util.Defer(func(err interface{}, stack string) {
				wg.Done()
				if err != nil {
					logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err, "stack": stack}).Info("ping，异常")
				}
			})

			_, err := this.Ping(ctx, address)
			if err != nil {
				return
			}
			listChan <- address
		}(addresses[i])
	}
	wg.Wait()

	close(listChan)
	list := make([]string, 0, len(listChan))
	for address := range listChan {
		list = append(list, address)
	}
	return list
}
func (this *ServerCenterClient) Ping(ctx context.Context, address string) (*common_model.PingResponse, error) {
	ctx = util.RmReqId(ctx)
	type Response struct {
		common_model.HttpResponse
		Data common_model.PingResponse `json:"data"`
	}
	var response Response
	err := util.HttpApiWithTry(ctx, "ping", this.try, nil, &response, func() (*resty.Response, error) {
		response, err := this.httpClient.R().SetContext(ctx).
			SetHeader(this.genJWT(ctx)).
			Post(this.getUrl(ctx, address, common_model.PingPath))
		return response, err
	})
	return &response.Data, err
}

func (this *ServerCenterClient) GetUrl(ctx context.Context, path string) string {
	return this.getUrl(ctx, this.getAddress(ctx), path)
}
func (this *ServerCenterClient) getUrl(ctx context.Context, address, path string) string {
	if strings.HasSuffix(address, "/") && strings.HasPrefix(path, "/") && len(path) > 0 {
		path = path[1:]
	}
	return address + path
}
func (this *ServerCenterClient) getAddress(ctx context.Context) string {
	list := this.handler.ListAddress(ctx)
	if len(list) == 0 {
		return ""
	}
	logId := util.GetLogId(ctx)
	index := int(logId) % len(list)
	return list[index]
}
func (this *ServerCenterClient) genJWT(ctx context.Context) (string, string) {
	return util.GenAuthorizationJWT(ctx, this.timeout, this.handler.GetSecret(ctx))
}
func (this *ServerCenterClient) getName(ctx context.Context) string {
	return this.name
}
