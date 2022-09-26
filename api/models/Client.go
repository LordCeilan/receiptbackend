package models

type Client struct {
	ID       uint32 `gorm:"primary_key;auto_increment" json:"id"`
	NickName string `gorm:"size:255;not null; unique" json:"nickname"`
	Email    string `gorm:"size:100;not null;unique" json:"email"`
}
