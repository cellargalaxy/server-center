package model

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/sirupsen/logrus"
)

const (
	ListenAddress         = ":7557"
	AddServerConfPath     = "/api/addServerConf"
	RemoveServerConfPath  = "/api/removeServerConf"
	GetLastServerConfPath = "/api/getLastServerConf"
	ListServerConfPath    = "/api/listServerConf"
	ListAllServerNamePath = "/api/listAllServerName"
)

type Config struct {
	LogLevel logrus.Level `yaml:"log_level" json:"log_level"`
	MysqlDsn string       `yaml:"mysql_dsn" json:"-"`
	ShowSql  bool         `yaml:"show_sql" json:"show_sql"`
	Secret   string       `yaml:"secret" json:"-"`

	PullSyncCron   string `yaml:"pull_sync_cron" json:"pull_sync_cron"`
	PullSyncHost   string `yaml:"pull_sync_host" json:"pull_sync_host"`
	PullSyncSecret string `yaml:"pull_sync_secret" json:"-"`
}

func (this Config) String() string {
	return util.ToJsonString(this)
}
