package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Client struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null; unique" json:"name"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Receipts  []Receipt `gorm:"foreignKey:UserRefer"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *Client) Prepare() {
	c.ID = 0
	c.Name = html.EscapeString(strings.TrimSpace(c.Name))
	c.Email = html.EscapeString(strings.TrimSpace(c.Email))
	c.Receipts = []Receipt{}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Client) Validate() error {
	if c.Name == "" {
		return errors.New("required Tittle")
	}

	if c.Email == "" {
		return errors.New("required Email")
	}

	return nil
}

func (c *Client) SaveClient(db *gorm.DB) (*Client, error) {
	// var err error
	err := db.Debug().Model(&Client{}).Create(&c).Error
	if err != nil {
		return &Client{}, err
	}

	// if c.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", c.ReceiptId).Take(&c.)
	// }
	return c, nil
}
