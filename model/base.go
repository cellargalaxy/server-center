package model

import "time"

type Model struct {
	Id        int       `gorm:"id" json:"id"`
	CreatedAt time.Time `gorm:"created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at" json:"updated_at"`
}