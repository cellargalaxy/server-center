package model

import "github.com/cellargalaxy/go_common/util"

type ServerConf struct {
	ServerName string `json:"server_name" gorm:"data"`
	Version    int    `json:"version" gorm:"version"`
	Remark     string `json:"remark" gorm:"remark"`
	ConfText   string `json:"conf_text" gorm:"conf_text"`
}

func (this ServerConf) String() string {
	return util.ToJsonString(this)
}

type ServerConfModel struct {
	Model
	ServerConf
}

func (this ServerConfModel) String() string {
	return util.ToJsonString(this)
}

func (ServerConfModel) TableName() string {
	return "server_conf"
}

type ServerConfInquiry struct {
	ServerConfModel
	CurrentVersion int `json:"version"`
}

func (this ServerConfInquiry) String() string {
	return util.ToJsonString(this)
}
