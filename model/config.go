package model

import (
	"github.com/cellargalaxy/go_common/util"
)

const (
	DefaultServerName     = "server_center"
	ListenAddress         = ":7557"
	StaticPath            = "/static"
	AddServerConfPath     = "/api/addServerConf"
	RemoveServerConfPath  = "/api/removeServerConf"
	GetLastServerConfPath = "/api/getLastServerConf"
	ListServerConfPath    = "/api/listServerConf"
	ListAllServerNamePath = "/api/listAllServerName"
	AddEventPath          = "/api/addEvent"
)

type Config struct {
	MysqlDsn  string   `yaml:"mysql_dsn" json:"-"`
	ShowSql   bool     `yaml:"show_sql" json:"show_sql"`
	Addresses []string `yaml:"addresses" json:"addresses"`
	Secret    string   `yaml:"secret" json:"-"`

	PullSyncCron   string `yaml:"pull_sync_cron" json:"pull_sync_cron"`
	PullSyncHost   string `yaml:"pull_sync_host" json:"pull_sync_host"`
	PullSyncSecret string `yaml:"pull_sync_secret" json:"-"`

	EventMaxSave   int    `yaml:"event_max_save" json:"event_max_save"`
	ClearEventCron string `yaml:"clear_event_cron" json:"clear_event_cron"`
}

func (this Config) String() string {
	return util.ToJsonString(this)
}
