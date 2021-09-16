package model

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/sirupsen/logrus"
)

const (
	ListenAddress = ":7557"
)

type Config struct {
	LogLevel logrus.Level `ini:"log_level" json:"log_level"`
	MysqlDsn string       `ini:"mysql_dsn" json:"-"`
	ShowSql  bool         `ini:"show_sql" json:"show_sql"`
	Secret   string       `ini:"secret" json:"-"`
}

func (this Config) String() string {
	return util.ToJsonString(this)
}
