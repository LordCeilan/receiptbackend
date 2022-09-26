package models

import "time"

type Receipt struct {
	ClientID  uint64    `gorm:"size:255;not null:unique" json:"clientid"`
	Contents  string    `gorm:"size:255;not null:unique" json:"content"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
