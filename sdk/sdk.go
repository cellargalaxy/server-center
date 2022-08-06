package sdk

import (
	"context"
	"fmt"
	common_model "github.com/cellargalaxy/go_common/model"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/go-resty/resty/v2"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"net/http"
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

type ServerCenterClient struct {
	timeout    time.Duration
	retry      int
	httpClient *resty.Client
	handler    ServerCenterHandlerInter

	conf model.ServerConfModel
	lock sync.Mutex
	pool *ants.Pool
}

func NewDefaultServerCenterClient(ctx context.Context, handler ServerCenterHandlerInter) (*ServerCenterClient, error) {
	var err error
	client, err = NewServerCenterClient(ctx, util.TimeoutDefault, util.RetryDefault, util.GetHttpClient(), handler)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewServerCenterClient(ctx context.Context, timeout time.Duration, retry int, httpClient *resty.Client, handler ServerCenterHandlerInter) (*ServerCenterClient, error) {
	if handler == nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Error("创建ServerCenterClient，handler为空")
		return nil, fmt.Errorf("创建ServerCenterClient，handler为空")
	}
	client := &ServerCenterClient{timeout: timeout, retry: retry, httpClient: httpClient, handler: handler}
	return client, nil
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

	if this.pool != nil && !this.pool.IsClosed() {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"name": this.getName(ctx)}).Warn("已启动")
		return nil
	}

	var err error
	this.pool, err = util.NewForeverSingleGoPool(ctx, this.getName(ctx), time.Second, this.flushServerConf)
	if err != nil {
		if this.pool != nil {
			this.pool.Release()
		}
		this.pool = nil
		return err
	}

	return nil
}
func (this *ServerCenterClient) flushServerConf(ctx context.Context, cancel func()) {
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
	var jsonString string
	var object *model.ServerConfModel
	var err error
	for i := 0; i < this.retry; i++ {
		jsonString, err = this.requestGetLastServerConf(ctx, serverName)
		if err == nil {
			object, err = this.analysisGetLastServerConf(ctx, jsonString)
			if err == nil {
				return object, err
			}
		}
	}
	return object, err
}
func (this *ServerCenterClient) analysisGetLastServerConf(ctx context.Context, jsonString string) (*model.ServerConfModel, error) {
	type Response struct {
		Code int                             `json:"code"`
		Msg  string                          `json:"msg"`
		Data model.GetLastServerConfResponse `json:"data"`
	}
	var response Response
	err := util.UnmarshalJsonString(jsonString, &response)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("查询最新服务配置，解析响应异常")
		return nil, fmt.Errorf("查询最新服务配置，解析响应异常")
	}
	if response.Code != util.HttpSuccessCode {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"response": util.ToJsonString(response)}).Error("查询最新服务配置，失败")
		return nil, fmt.Errorf("查询最新服务配置，失败")
	}
	return response.Data.Conf, nil
}
func (this *ServerCenterClient) requestGetLastServerConf(ctx context.Context, serverName string) (string, error) {
	response, err := this.httpClient.R().SetContext(ctx).
		SetHeader(this.genJWT(ctx)).
		SetQueryParam("server_name", serverName).
		Get(this.GetUrl(ctx, model.GetLastServerConfPath))

	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("查询最新服务配置，请求异常")
		return "", fmt.Errorf("查询最新服务配置，请求异常")
	}
	if response == nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Error("查询最新服务配置，响应为空")
		return "", fmt.Errorf("查询最新服务配置，响应为空")
	}
	statusCode := response.StatusCode()
	body := response.String()
	logrus.WithContext(ctx).WithFields(logrus.Fields{"statusCode": statusCode, "len(body)": len(body)}).Info("查询最新服务配置，响应")
	if statusCode != http.StatusOK {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"StatusCode": statusCode}).Error("查询最新服务配置，响应码失败")
		return "", fmt.Errorf("查询最新服务配置，响应码失败: %+v", statusCode)
	}
	return body, nil
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
	var jsonString string
	var object []string
	var err error
	for i := 0; i < this.retry; i++ {
		jsonString, err = this.requestListAllServerName(ctx)
		if err == nil {
			object, err = this.analysisListAllServerName(ctx, jsonString)
			if err == nil {
				return object, err
			}
		}
	}
	return object, err
}
func (this *ServerCenterClient) analysisListAllServerName(ctx context.Context, jsonString string) ([]string, error) {
	type Response struct {
		Code int                             `json:"code"`
		Msg  string                          `json:"msg"`
		Data model.ListAllServerNameResponse `json:"data"`
	}
	var response Response
	err := util.UnmarshalJsonString(jsonString, &response)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("查询服务配置列表，解析响应异常")
		return nil, fmt.Errorf("查询服务配置列表，解析响应异常")
	}
	if response.Code != util.HttpSuccessCode {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"response": util.ToJsonString(response)}).Error("查询服务配置列表，失败")
		return nil, fmt.Errorf("查询服务配置列表，失败")
	}
	return response.Data.List, nil
}
func (this *ServerCenterClient) requestListAllServerName(ctx context.Context) (string, error) {
	response, err := this.httpClient.R().SetContext(ctx).
		SetHeader(this.genJWT(ctx)).
		Get(this.GetUrl(ctx, model.ListAllServerNamePath))

	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("查询服务配置列表，请求异常")
		return "", fmt.Errorf("查询服务配置列表，请求异常")
	}
	if response == nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Error("查询服务配置列表，响应为空")
		return "", fmt.Errorf("查询服务配置列表，响应为空")
	}
	statusCode := response.StatusCode()
	body := response.String()
	logrus.WithContext(ctx).WithFields(logrus.Fields{"statusCode": statusCode, "len(body)": len(body)}).Info("查询服务配置列表，响应")
	if statusCode != http.StatusOK {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"StatusCode": statusCode}).Error("查询服务配置列表，响应码失败")
		return "", fmt.Errorf("查询服务配置列表，响应码失败: %+v", statusCode)
	}
	return body, nil
}

func (this *ServerCenterClient) AddEvent(ctx context.Context, object []model.Event) error {
	ctx = util.SetReqId(ctx)
	var jsonString string
	var err error
	for i := 0; i < this.retry; i++ {
		jsonString, err = this.requestAddEvent(ctx, object)
		if err == nil {
			err = this.analysisAddEvent(ctx, jsonString)
			if err == nil {
				return err
			}
		}
	}
	return err
}
func (this *ServerCenterClient) analysisAddEvent(ctx context.Context, jsonString string) error {
	type Response struct {
		Code int                    `json:"code"`
		Msg  string                 `json:"msg"`
		Data model.AddEventResponse `json:"data"`
	}
	var response Response
	err := util.UnmarshalJsonString(jsonString, &response)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("插入批量事件，解析响应异常")
		return fmt.Errorf("插入批量事件，解析响应异常")
	}
	if response.Code != util.HttpSuccessCode {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"response": util.ToJsonString(response)}).Error("插入批量事件，失败")
		return fmt.Errorf("插入批量事件，失败")
	}
	return nil
}
func (this *ServerCenterClient) requestAddEvent(ctx context.Context, object []model.Event) (string, error) {
	var req model.AddEventRequest
	req.List = object
	response, err := this.httpClient.R().SetContext(ctx).
		SetHeader(this.genJWT(ctx)).
		SetBody(req).
		Post(this.GetUrl(ctx, model.AddEventPath))

	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("插入批量事件，请求异常")
		return "", fmt.Errorf("插入批量事件，请求异常")
	}
	if response == nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Error("插入批量事件，响应为空")
		return "", fmt.Errorf("插入批量事件，响应为空")
	}
	statusCode := response.StatusCode()
	body := response.String()
	logrus.WithContext(ctx).WithFields(logrus.Fields{"statusCode": statusCode, "body": body}).Info("插入批量事件，响应")
	if statusCode != http.StatusOK {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"StatusCode": statusCode}).Error("插入批量事件，响应码失败")
		return "", fmt.Errorf("插入批量事件，响应码失败: %+v", statusCode)
	}
	return body, nil
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
	var jsonString string
	var object *common_model.PingResponse
	var err error
	for i := 0; i < this.retry; i++ {
		jsonString, err = this.requestPing(ctx, address)
		if err == nil {
			object, err = this.analysisPing(ctx, jsonString)
			if err == nil {
				return object, err
			}
		}
	}
	return object, err
}
func (this *ServerCenterClient) analysisPing(ctx context.Context, jsonString string) (*common_model.PingResponse, error) {
	type Response struct {
		Code int                       `json:"code"`
		Msg  string                    `json:"msg"`
		Data common_model.PingResponse `json:"data"`
	}
	var response Response
	err := util.UnmarshalJsonString(jsonString, &response)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("ping，解析响应异常")
		return nil, fmt.Errorf("ping，解析响应异常")
	}
	if response.Code != util.HttpSuccessCode {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"response": util.ToJsonString(response)}).Error("ping，失败")
		return nil, fmt.Errorf("ping，失败")
	}
	return &response.Data, nil
}
func (this *ServerCenterClient) requestPing(ctx context.Context, address string) (string, error) {
	response, err := this.httpClient.R().SetContext(ctx).
		SetHeader(this.genJWT(ctx)).
		Post(this.getUrl(ctx, address, util.PingPath))

	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("ping，请求异常")
		return "", fmt.Errorf("ping，请求异常")
	}
	if response == nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Error("ping，响应为空")
		return "", fmt.Errorf("ping，响应为空")
	}
	statusCode := response.StatusCode()
	body := response.String()
	logrus.WithContext(ctx).WithFields(logrus.Fields{"statusCode": statusCode, "body": body}).Info("ping，响应")
	if statusCode != http.StatusOK {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"StatusCode": statusCode}).Error("ping，响应码失败")
		return "", fmt.Errorf("ping，响应码失败: %+v", statusCode)
	}
	return body, nil
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
	if this.handler == nil {
		return "ServerCenterClient"
	}
	return fmt.Sprintf("ServerCenterClient-%s", this.handler.GetServerName(ctx))
}
