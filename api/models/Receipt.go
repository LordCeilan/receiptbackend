package models

import "time"

type Receipt struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Client    string    `gorm:"size:255;not null:unique" json:"client"`
	Content   string    `gorm:"size:255;not null:unique" json:"content"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
