package model

import (
	"github.com/cellargalaxy/go_common/util"
	"time"
)

type Event struct {
	LogId      int64     `json:"log_id" form:"log_id" query:"log_id" gorm:"log_id;not null;default:0;index:idx_log_id"`
	ServerName string    `json:"server_name" form:"server_name" query:"server_name" gorm:"server_name;not null;default:'';index:idx_server_name"`
	Ip         string    `json:"ip" form:"ip" query:"ip" gorm:"ip;not null;default:''"`
	Group      string    `json:"group" form:"group" query:"group" gorm:"group;not null;default:'';index:idx_group"`
	Name       string    `json:"name" form:"name" query:"name" gorm:"name;not null;default:'';index:idx_name"`
	Value      float64   `json:"value" form:"value" query:"value" gorm:"value;not null;default:0"`
	Data       string    `json:"data" form:"data" query:"data" gorm:"data;not null;default:''"`
	CreateTime time.Time `json:"create_time" form:"create_time" query:"create_time" gorm:"create_time"`
}

func (this Event) String() string {
	return util.ToJsonString(this)
}

type EventModel struct {
	Id int `json:"id" form:"id" query:"id" gorm:"id;auto_increment;primary_key"`
	Event
}

func (this EventModel) String() string {
	return util.ToJsonString(this)
}

func (EventModel) TableName() string {
	return "event"
}

type EventInquiry struct {
	EventModel
	EndCreateTime time.Time `json:"end_create_time" form:"end_create_time" query:"end_create_time"`
	Offset        int       `json:"offset" form:"offset" query:"offset"`
	Limit         int       `json:"limit" form:"limit" query:"limit"`
}

func (this EventInquiry) String() string {
	return util.ToJsonString(this)
}

type AddEventRequest struct {
	List []Event
}

func (this AddEventRequest) String() string {
	return "-"
}

type AddEventResponse struct {
}

func (this AddEventResponse) String() string {
	return "-"
}

type RemoveEventRequest struct {
	EventInquiry
}

func (this RemoveEventRequest) String() string {
	return util.ToJsonString(this)
}

type RemoveEventResponse struct {
}

func (this RemoveEventResponse) String() string {
	return util.ToJsonString(this)
}
