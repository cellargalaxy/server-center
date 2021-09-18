package sdk

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/cellargalaxy/go_common/consd"
	common_model "github.com/cellargalaxy/go_common/model"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/model"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	serverCenterAddressEnvKey = "server_center_address"
	serverCenterSecretEnvKey  = "server_center_secret"
)

func GetEnvServerCenterAddress(ctx context.Context) string {
	return util.GetEnvString(serverCenterAddressEnvKey, "")
}
func GetEnvServerCenterSecret(ctx context.Context) string {
	return util.GetEnvString(serverCenterSecretEnvKey, "")
}

type ServerCenterHandlerInter interface {
	GetAddress(ctx context.Context) string
	GetSecret(ctx context.Context) string
	GetInterval(ctx context.Context) time.Duration
	ParseConf(ctx context.Context, object model.ServerConfModel) error
}

type ServerCenterClient struct {
	retry      int
	httpClient *resty.Client
	handler    ServerCenterHandlerInter
	conf       model.ServerConfModel
	running    bool
}

func NewDefaultServerCenterClient(handler ServerCenterHandlerInter) (*ServerCenterClient, error) {
	return NewServerCenterClient(3*time.Second, 3*time.Second, 3, handler)
}

func NewServerCenterClient(timeout, sleep time.Duration, retry int, handler ServerCenterHandlerInter) (*ServerCenterClient, error) {
	if handler == nil {
		return nil, fmt.Errorf("handler为空")
	}
	httpClient := createHttpClient(timeout, sleep, retry)
	return &ServerCenterClient{retry: retry, httpClient: httpClient, handler: handler}, nil
}

func createHttpClient(timeout, sleep time.Duration, retry int) *resty.Client {
	httpClient := resty.New().
		SetTimeout(timeout).
		SetRetryCount(retry).
		SetRetryWaitTime(sleep).
		SetRetryMaxWaitTime(5 * time.Minute).
		AddRetryCondition(func(response *resty.Response, err error) bool {
			ctx := util.CreateLogCtx()
			if response != nil && response.Request != nil {
				ctx = response.Request.Context()
			}
			var statusCode int
			if response != nil {
				statusCode = response.StatusCode()
			}
			isRetry := statusCode != http.StatusOK || err != nil
			if isRetry {
				logrus.WithContext(ctx).WithFields(logrus.Fields{"statusCode": statusCode, "err": err}).Warn("HTTP请求异常，进行重试")
			}
			return isRetry
		}).
		SetRetryAfter(func(client *resty.Client, response *resty.Response) (time.Duration, error) {
			ctx := util.CreateLogCtx()
			if response != nil && response.Request != nil {
				ctx = response.Request.Context()
			}
			var attempt int
			if response != nil && response.Request != nil {
				attempt = response.Request.Attempt
			}
			if attempt > retry {
				logrus.WithContext(ctx).WithFields(logrus.Fields{"attempt": attempt}).Error("HTTP请求异常，超过最大重试次数")
				return 0, fmt.Errorf("HTTP请求异常，超过最大重试次数")
			}
			duration := util.WareDuration(sleep)
			for i := 0; i < attempt-1; i++ {
				duration *= 10
			}
			logrus.WithContext(ctx).WithFields(logrus.Fields{"attempt": attempt, "duration": duration}).Warn("HTTP请求异常，休眠重试")
			return duration, nil
		}).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	return httpClient
}

func (this *ServerCenterClient) StartConfWithInitConf(ctx context.Context) (*model.ServerConfModel, error) {
	if this.running {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Warn("ServerCenterClient开始，已开始")
		return nil, nil
	}
	object, err := this.GetAndParseLastServerConf(ctx)
	if err != nil {
		return nil, err
	}
	this.running = true
	go this.startConf()
	return object, nil
}

func (this *ServerCenterClient) StartConf(ctx context.Context) {
	if this.running {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Warn("ServerCenterClient开始，已开始")
		return
	}
	this.running = true
	go this.startConf()
}

func (this *ServerCenterClient) startConf() {
	for true {
		ctx := util.CreateLogCtx()
		this.GetAndParseLastServerConf(ctx)
		time.Sleep(util.WareDuration(this.handler.GetInterval(ctx)))
	}
}

func (this *ServerCenterClient) GetAndParseLastServerConf(ctx context.Context) (*model.ServerConfModel, error) {
	object, err := this.GetLastServerConf(ctx)
	if err != nil {
		return nil, err
	}
	if object == nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"current_version": this.conf.Version}).Info("查询并解析最新服务配置，服务配置无更新")
		return nil, nil
	}
	if object.Version <= this.conf.Version {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"current_version": this.conf.Version, "last_version": object.Version}).Warn("查询并解析最新服务配置，服务配置无更新")
		return nil, nil
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{"current_version": this.conf.Version, "last_version": object.Version}).Info("查询并解析最新服务配置，服务配置更新")
	err = this.handler.ParseConf(ctx, *object)
	if err == nil {
		this.conf = *object
	}
	return object, err
}

func (this *ServerCenterClient) GetLastServerConf(ctx context.Context) (*model.ServerConfModel, error) {
	var jwtToken, jsonString string
	var object *model.ServerConfModel
	var err error
	for i := 0; i < this.retry; i++ {
		jwtToken, err = this.genJWT(ctx)
		if err != nil {
			return nil, err
		}
		jsonString, err = this.requestGetLastServerConf(ctx, jwtToken)
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
	if response.Code != consd.HttpSuccessCode {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"code": response.Code, "msg": response.Msg}).Error("查询最新服务配置，失败")
		return nil, fmt.Errorf("查询最新服务配置，失败")
	}
	return response.Data.Conf, nil
}

func (this *ServerCenterClient) requestGetLastServerConf(ctx context.Context, jwtToken string) (string, error) {
	response, err := this.httpClient.R().SetContext(ctx).
		SetHeader(util.LogIdKey, fmt.Sprint(util.GetLogId(ctx))).
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetQueryParam("server_name", util.GetServerNameWithPanic()).
		SetQueryParam("current_version", fmt.Sprint(this.conf.Version)).
		Get(this.handler.GetAddress(ctx) + "/api/getLastServerConf")

	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("查询最新服务配置，请求异常")
		return "", fmt.Errorf("查询最新服务配置，请求异常")
	}
	if response == nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("查询最新服务配置，响应为空")
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

func (this *ServerCenterClient) genJWT(ctx context.Context) (string, error) {
	now := time.Now()
	var claims common_model.Claims
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = now.Unix() + int64(this.retry*3)
	claims.RequestId = fmt.Sprint(util.GenId())
	jwtToken, err := util.GenJWT(ctx, this.handler.GetSecret(ctx), claims)
	return jwtToken, err
}
