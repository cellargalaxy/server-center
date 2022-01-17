package model

import "github.com/cellargalaxy/go_common/util"

type ServerConf struct {
	ServerName string `json:"server_name" gorm:"server_name" form:"server_name" query:"server_name"`
	Version    int    `json:"version" gorm:"version" form:"version" query:"version"`
	Remark     string `json:"remark" gorm:"remark" form:"remark" query:"remark"`
	ConfText   string `json:"conf_text" gorm:"conf_text" form:"conf_text" query:"conf_text"`
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
	CurrentVersion int `json:"current_version" form:"current_version" query:"current_version"`
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

type RemoveServerConfRequest struct {
	ServerConfInquiry
}

func (this RemoveServerConfRequest) String() string {
	return util.ToJsonString(this)
}

type RemoveServerConfResponse struct {
	Conf *ServerConfModel `json:"conf"`
}

func (this RemoveServerConfResponse) String() string {
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
