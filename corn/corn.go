package corn

import (
	"context"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/config"
	"github.com/cellargalaxy/server_center/service/service"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func Init(ctx context.Context) {
	cronObject := cron.New()

	if config.Config.PullSyncCron != "" {
		var job pullSyncJob
		job.Address = config.Config.PullSyncHost
		job.Secret = config.Config.PullSyncSecret
		entryId, err := cronObject.AddJob(config.Config.PullSyncCron, &job)
		if err != nil {
			panic(err)
		}
		logrus.WithContext(ctx).WithFields(logrus.Fields{"pullSyncJob": job, "entryId": entryId}).Info("定时任务，添加定时")
	}

	if config.Config.ClearEventCron != "" {
		var job clearEventJob
		entryId, err := cronObject.AddJob(config.Config.ClearEventCron, &job)
		if err != nil {
			panic(err)
		}
		logrus.WithContext(ctx).WithFields(logrus.Fields{"clearEventJob": job, "entryId": entryId}).Info("定时任务，添加定时")
	}

	if config.Config.ClearConfigCron != "" {
		var job clearConfigJob
		entryId, err := cronObject.AddJob(config.Config.ClearConfigCron, &job)
		if err != nil {
			panic(err)
		}
		logrus.WithContext(ctx).WithFields(logrus.Fields{"clearConfigJob": job, "entryId": entryId}).Info("定时任务，添加定时")
	}

	cronObject.Start()
	logrus.WithContext(ctx).WithFields(logrus.Fields{}).Info("定时任务，添加完成")
}

type pullSyncJob struct {
	Address string `json:"address"`
	Secret  string `json:"-"`
}

func (this pullSyncJob) String() string {
	return util.ToJsonString(this)
}

func (this *pullSyncJob) Run() {
	ctx := util.GenCtx()
	logrus.WithContext(ctx).WithFields(logrus.Fields{"pullSyncJob": this}).Info("定时任务，执行任务开始")
	service.PullSync(ctx, this.Address, this.Secret)
	logrus.WithContext(ctx).WithFields(logrus.Fields{"pullSyncJob": this}).Info("定时任务，执行任务完成")
}

type clearEventJob struct {
}

func (this clearEventJob) String() string {
	return util.ToJsonString(this)
}

func (this *clearEventJob) Run() {
	ctx := util.GenCtx()
	logrus.WithContext(ctx).WithFields(logrus.Fields{"clearEventJob": this}).Info("定时任务，执行任务开始")
	service.ClearEvent(ctx)
	logrus.WithContext(ctx).WithFields(logrus.Fields{"clearEventJob": this}).Info("定时任务，执行任务完成")
}

type clearConfigJob struct {
}

func (this clearConfigJob) String() string {
	return util.ToJsonString(this)
}

func (this *clearConfigJob) Run() {
	ctx := util.GenCtx()
	logrus.WithContext(ctx).WithFields(logrus.Fields{"clearConfigJob": this}).Info("定时任务，执行任务开始")
	service.ClearConfig(ctx)
	logrus.WithContext(ctx).WithFields(logrus.Fields{"clearConfigJob": this}).Info("定时任务，执行任务完成")
}
