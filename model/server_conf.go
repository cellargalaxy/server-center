package model

import "github.com/cellargalaxy/go_common/util"

type ServerConf struct {
	ServerName string `json:"server_name" gorm:"data" form:"data"`
	Version    int    `json:"version" gorm:"version" form:"version"`
	Remark     string `json:"remark" gorm:"remark" form:"remark"`
	ConfText   string `json:"conf_text" gorm:"conf_text" form:"conf_text"`
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
	CurrentVersion int `json:"current_version" form:"current_version"`
}

func (this ServerConfInquiry) String() string {
	return util.ToJsonString(this)
}

type AddServerConfRequest struct {
	ServerConf
}

func (this AddServerConfRequest) String() string {
	return util.ToJsonString(this)
}

type AddServerConfResponse struct {
	Conf *ServerConfModel `json:"conf"`
}

func (this AddServerConfResponse) String() string {
	return util.ToJsonString(this)
}

type GetLastServerConfRequest struct {
	ServerConfInquiry
}

func (this GetLastServerConfRequest) String() string {
	return util.ToJsonString(this)
}

type GetLastServerConfResponse struct {
	Conf *ServerConfModel `json:"conf"`
}

func (this GetLastServerConfResponse) String() string {
	return util.ToJsonString(this)
}

type ListServerConfRequest struct {
	ServerConfInquiry
}

func (this ListServerConfRequest) String() string {
	return util.ToJsonString(this)
}

type ListServerConfResponse struct {
	List []ServerConfModel `json:"list"`
}

func (this ListServerConfResponse) String() string {
	return util.ToJsonString(this)
}

type ListAllServerNameRequest struct {
}

func (this ListAllServerNameRequest) String() string {
	return util.ToJsonString(this)
}

type ListAllServerNameResponse struct {
	List []string `json:"list"`
}

func (this ListAllServerNameResponse) String() string {
	return util.ToJsonString(this)
}
