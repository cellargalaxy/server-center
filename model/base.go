package model

import "time"

type Model struct {
	Id        int       `gorm:"id" json:"id" form:"id" query:"id"`
	CreatedAt time.Time `gorm:"created_at" json:"created_at" form:"created_at" query:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at" json:"updated_at" form:"updated_at" query:"updated_at"`
}
