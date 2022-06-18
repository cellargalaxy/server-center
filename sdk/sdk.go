package sdk

import (
	"context"
	"fmt"
	common_model "github.com/cellargalaxy/go_common/model"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	serverCenterAddressEnvKey = "server_center_address"
	serverCenterSecretEnvKey  = "server_center_secret"
)

var addresses []string
var secret string
var serverCenterClient *ServerCenterClient

func init() {
	ctx := util.CreateLogCtx()
	var err error

	address := GetEnvServerCenterAddress(ctx)
	if address != "" {
		addresses = append(addresses, address)
	}
	secret = GetEnvServerCenterSecret(ctx)

	var handler ServerCenterHandler
	serverCenterClient, err = NewDefaultServerCenterClient(ctx, &handler)
	if err != nil {
		panic(err)
	}
	if serverCenterClient == nil {
		panic("创建serverCenterClient为空")
	}
	serverCenterClient.StartConfWithInitConf(ctx)
}

func GetEnvServerCenterAddress(ctx context.Context) string {
	return util.GetEnvString(serverCenterAddressEnvKey, "")
}
func GetEnvServerCenterSecret(ctx context.Context) string {
	return util.GetEnvString(serverCenterSecretEnvKey, "")
}
func GetEnvServerName(ctx context.Context, defaultServerName string) string {
	return util.GetServerName(defaultServerName)
}
func ListAddress(ctx context.Context) []string {
	return addresses
}

type ServerCenterHandlerInter interface {
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
	conf       model.ServerConfModel
	once       sync.Once
}

func NewDefaultServerCenterClient(ctx context.Context, handler ServerCenterHandlerInter) (*ServerCenterClient, error) {
	return NewServerCenterClient(ctx, 3*time.Second, 3, handler)
}

func NewServerCenterClient(ctx context.Context, timeout time.Duration, retry int, handler ServerCenterHandlerInter) (*ServerCenterClient, error) {
	if handler == nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Error("创建ServerCenterClient，handler为空")
		return nil, fmt.Errorf("创建ServerCenterClient，handler为空")
	}
	if handler.GetServerName(ctx) == "" {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Error("创建ServerCenterClient，ServerName为空")
		return nil, fmt.Errorf("创建ServerCenterClient，ServerName为空")
	}
	httpClient := util.CreateNotTryHttpClient(timeout)
	client := &ServerCenterClient{timeout: timeout, retry: retry, httpClient: httpClient, handler: handler}
	client.conf.Version = math.MinInt32
	return client, nil
}

func (this *ServerCenterClient) ResetVersion(ctx context.Context) {
	this.conf.Version = 0
}

func (this *ServerCenterClient) StartConfWithInitConf(ctx context.Context) {
	for {
		_, err := this.GetAndParseLastServerConf(ctx)
		if err == nil {
			break
		}
		time.Sleep(util.WareDuration(time.Second))
	}
	this.StartServerCenter(ctx)
}
func (this *ServerCenterClient) StartServerCenter(ctx context.Context) {
	this.once.Do(this.startServerCenterAsync)
}
func (this *ServerCenterClient) startServerCenterAsync() {
	this.startConfAsync()
}
func (this *ServerCenterClient) startConfAsync() {
	go func() {
		ctx := util.CreateLogCtx()
		defer util.Defer(ctx, func(ctx context.Context, err interface{}, stack string) {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err, "stack": stack}).Warn("startConfAsync，结束")
			time.Sleep(util.WareDuration(util.MaxDuration(this.handler.GetInterval(ctx), time.Minute*5)))
			this.startConfAsync()
		})

		for {
			ctx := util.CreateLogCtx()
			this.GetAndParseLastServerConf(ctx)
			time.Sleep(util.WareDuration(this.handler.GetInterval(ctx)))
		}
	}()
}
func (this *ServerCenterClient) GetAndParseLastServerConf(ctx context.Context) (*model.ServerConfModel, error) {
	object, err := this.GetLastServerConf(ctx)
	if err != nil {
		return nil, err
	}
	if object == nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"server_name": this.handler.GetServerName(ctx), "current_version": this.conf.Version}).Info("查询并解析最新服务配置，服务配置无更新")
		return nil, nil
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{"server_name": this.handler.GetServerName(ctx), "current_version": this.conf.Version, "last_version": object.Version}).Info("查询并解析最新服务配置，查询服务配置")
	if object.Version <= this.conf.Version {
		return nil, nil
	}
	err = this.handler.ParseConf(ctx, *object)
	if err != nil {
		return object, err
	}
	this.conf = *object
	this.saveLocalFileServerConf(ctx, this.conf.ConfText)
	return object, nil
}

func (this *ServerCenterClient) GetLastServerConf(ctx context.Context) (*model.ServerConfModel, error) {
	if len(addresses) == 0 {
		return this.GetLocalFileServerConf(ctx)
	}
	object, err := this.GetRemoteLastServerConf(ctx)
	if err == nil {
		return object, nil
	}
	return this.GetLocalFileServerConf(ctx)
}
func (this *ServerCenterClient) saveLocalFileServerConf(ctx context.Context, confText string) error {
	filePath, err := this.getLocalFilePath(ctx)
	if err != nil {
		return err
	}
	return util.WriteFileWithString(ctx, filePath, confText)
}
func (this *ServerCenterClient) GetLocalFileServerConf(ctx context.Context) (*model.ServerConfModel, error) {
	filePath, err := this.getLocalFilePath(ctx)
	if err != nil {
		return nil, err
	}
	confText, err := util.ReadFileWithString(ctx, filePath, "")
	if err != nil {
		return nil, err
	}
	if confText == "" {
		confText = this.handler.GetDefaultConf(ctx)
		if confText != "" {
			util.WriteFileWithString(ctx, filePath, confText)
		}
	}
	var serverConf model.ServerConfModel
	serverConf.ServerName = this.handler.GetServerName(ctx)
	serverConf.Version = this.conf.Version + 1 //本地文件更新了也能更新配置
	serverConf.ConfText = confText
	return &serverConf, nil
}
func (this *ServerCenterClient) getLocalFilePath(ctx context.Context) (string, error) {
	filePath := "resource/" + this.handler.GetServerName(ctx) + ".yml"
	logrus.WithContext(ctx).WithFields(logrus.Fields{"filePath": filePath}).Info("查询本地文件服务配置")
	return filePath, nil
}
func (this *ServerCenterClient) GetRemoteLastServerConf(ctx context.Context) (*model.ServerConfModel, error) {
	return this.GetRemoteLastServerConfByServerName(ctx, this.handler.GetServerName(ctx))
}
func (this *ServerCenterClient) GetRemoteLastServerConfByServerName(ctx context.Context, serverName string) (*model.ServerConfModel, error) {
	var jwtToken, jsonString string
	var object *model.ServerConfModel
	var err error
	for i := 0; i < this.retry; i++ {
		jwtToken, err = this.genJWT(ctx)
		if err != nil {
			return nil, err
		}
		jsonString, err = this.requestGetLastServerConf(ctx, serverName, jwtToken)
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
func (this *ServerCenterClient) requestGetLastServerConf(ctx context.Context, serverName, jwtToken string) (string, error) {
	response, err := this.httpClient.R().SetContext(ctx).
		SetHeader(util.LogIdKey, util.GetLogIdString(ctx)).
		SetHeader(util.GenAuthorizationHeader(ctx, jwtToken)).
		SetQueryParam("server_name", serverName).
		SetQueryParam("current_version", strconv.Itoa(this.conf.Version)).
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
	if len(addresses) == 0 {
		return this.ListLocalAllServerName(ctx)
	}
	object, err := this.ListRemoteAllServerName(ctx)
	if len(object) > 0 && err == nil {
		return object, nil
	}
	return this.ListLocalAllServerName(ctx)
}
func (this *ServerCenterClient) ListLocalAllServerName(ctx context.Context) ([]string, error) {
	return []string{this.handler.GetServerName(ctx)}, nil
}
func (this *ServerCenterClient) ListRemoteAllServerName(ctx context.Context) ([]string, error) {
	var jwtToken, jsonString string
	var object []string
	var err error
	for i := 0; i < this.retry; i++ {
		jwtToken, err = this.genJWT(ctx)
		if err != nil {
			return nil, err
		}
		jsonString, err = this.requestListAllServerName(ctx, jwtToken)
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
func (this *ServerCenterClient) requestListAllServerName(ctx context.Context, jwtToken string) (string, error) {
	response, err := this.httpClient.R().SetContext(ctx).
		SetHeader(util.LogIdKey, util.GetLogIdString(ctx)).
		SetHeader(util.GenAuthorizationHeader(ctx, jwtToken)).
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

func (this *ServerCenterClient) Ping(ctx context.Context, address string) (*common_model.PingResponse, error) {
	var jwtToken, jsonString string
	var object *common_model.PingResponse
	var err error
	for i := 0; i < this.retry; i++ {
		jwtToken, err = this.genJWT(ctx)
		if err != nil {
			return nil, err
		}
		jsonString, err = this.requestPing(ctx, jwtToken, address)
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
func (this *ServerCenterClient) requestPing(ctx context.Context, jwtToken, address string) (string, error) {
	response, err := this.httpClient.R().SetContext(ctx).
		SetHeader(util.LogIdKey, util.GetLogIdString(ctx)).
		SetHeader(util.GenAuthorizationHeader(ctx, jwtToken)).
		Get(this.getUrl(ctx, address, model.ListAllServerNamePath))

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
	list := addresses
	if len(list) == 0 {
		return ""
	}
	index := rand.Intn(len(list))
	return this.getUrl(ctx, list[index], path)
}
func (this *ServerCenterClient) getUrl(ctx context.Context, address, path string) string {
	if strings.HasSuffix(address, "/") && strings.HasPrefix(path, "/") && len(path) > 0 {
		path = path[1:]
	}
	return address + path
}

func (this *ServerCenterClient) genJWT(ctx context.Context) (string, error) {
	return util.GenDefaultJWT(ctx, this.timeout*time.Duration(this.retry+1), secret)
}
