package model

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/sirupsen/logrus"
)

const (
	ListenAddress = ":7557"
)

type Config struct {
	LogLevel logrus.Level `yaml:"log_level" json:"log_level"`
	MysqlDsn string       `yaml:"mysql_dsn" json:"-"`
	ShowSql  bool         `yaml:"show_sql" json:"show_sql"`
	Secret   string       `yaml:"secret" json:"-"`
}

func (this Config) String() string {
	return util.ToJsonString(this)
}
