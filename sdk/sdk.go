package sdk

import (
	"context"
	"crypto/tls"
	"encoding/json"
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

type ServerCenterHandlerInter interface {
	ParseConf(ctx context.Context, object model.ServerConfModel) error
}

type ServerCenterClient struct {
	address    string
	secret     string
	retry      int
	interval   time.Duration
	httpClient *resty.Client
	conf       model.ServerConfModel
	handler    ServerCenterHandlerInter
	running    bool
}

func NewDefaultServerCenterClient(handler ServerCenterHandlerInter) (*ServerCenterClient, error) {
	address := util.GetEnvString("server_center_address", "")
	secret := util.GetEnvString("server_center_secret", "")
	return NewServerCenterClient(3*time.Second, 3*time.Second, 3, address, secret, 5*time.Minute, handler)
}

func NewServerCenterClient(timeout, sleep time.Duration, retry int, address, secret string, interval time.Duration, handler ServerCenterHandlerInter) (*ServerCenterClient, error) {
	if address == "" {
		return nil, fmt.Errorf("address为空")
	}
	if secret == "" {
		return nil, fmt.Errorf("secret为空")
	}
	httpClient := createHttpClient(timeout, sleep, retry)
	return &ServerCenterClient{address: address, secret: secret, retry: retry, interval: interval, httpClient: httpClient, handler: handler}, nil
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
			retry := statusCode != http.StatusOK || err != nil
			if retry {
				logrus.WithContext(ctx).WithFields(logrus.Fields{"statusCode": statusCode, "err": err}).Warn("HTTP请求异常，进行重试")
			}
			return retry
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

func (this *ServerCenterClient) Start(ctx context.Context) (*model.ServerConfModel, error) {
	if this.running {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Warn("ServerCenterClient开始，已开始")
		return nil, nil
	}
	if this.handler == nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Warn("ServerCenterClient开始，ServerCenterHandlerInter为空")
		return nil, nil
	}
	object, err := this.GetAndParseLastServerConf(ctx)
	if err != nil {
		return nil, err
	}
	this.running = true
	go this.start()
	return object, nil
}

func (this *ServerCenterClient) start() {
	for true {
		ctx := util.CreateLogCtx()
		this.GetAndParseLastServerConf(ctx)
		time.Sleep(util.WareDuration(this.interval))
	}
}

func (this *ServerCenterClient) GetAndParseLastServerConf(ctx context.Context) (*model.ServerConfModel, error) {
	if this.handler == nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{}).Warn("查询并解析最新服务配置，ServerCenterHandlerInter为空")
		return nil, nil
	}
	object, err := this.GetLastServerConf(ctx)
	if object == nil || err != nil {
		return nil, err
	}
	err = this.handler.ParseConf(ctx, *object)
	if err == nil {
		this.conf = *object
	}
	return object, err
}

func (this *ServerCenterClient) GetLastServerConf(ctx context.Context) (*model.ServerConfModel, error) {
	var jsonString string
	var object *model.ServerConfModel
	var err error
	for i := 0; i < this.retry; i++ {
		jwtToken, err := this.genJWT(ctx)
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
	err := json.Unmarshal([]byte(jsonString), &response)
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
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetHeader(util.LogIdKey, fmt.Sprint(util.GetLogId(ctx))).
		SetQueryParam("server_name", util.GetServerNameWithPanic()).
		SetQueryParam("current_version", fmt.Sprint(this.conf.Version)).
		Get(this.address + "/api/getLastServerConf")

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
	jwtToken, err := util.GenJWT(ctx, this.secret, claims)
	return jwtToken, err
}
