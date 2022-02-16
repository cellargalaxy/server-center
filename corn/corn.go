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
		logrus.WithContext(ctx).WithFields(logrus.Fields{"corn": job, "entryId": entryId}).Info("定时任务，添加定时")
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
	ctx := util.CreateLogCtx()
	logrus.WithContext(ctx).WithFields(logrus.Fields{"corn": this}).Info("定时任务，执行任务开完")
	service.PullSync(ctx, this.Address, this.Secret)
	logrus.WithContext(ctx).WithFields(logrus.Fields{"corn": this}).Info("定时任务，执行任务完成")
}
