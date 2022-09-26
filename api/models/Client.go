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
	Receipts  Receipt   `gorm:"foreignKey:UserRefer json:receipt"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *Client) Prepare() {
	c.ID = 0
	c.Name = html.EscapeString(strings.TrimSpace(c.Name))
	c.Email = html.EscapeString(strings.TrimSpace(c.Email))
	c.Receipts = Receipt{}
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
	err := db.Debug().Create(&c).Error
	if err != nil {
		return &Client{}, err
	}
	return c, nil
}

func (c *Client) FindAllClients(db *gorm.DB) (*[]Client, error) {
	// var err error
	clients := []Client{}
	err := db.Debug().Model(&Client{}).Limit(100).Find(&clients).Error
	if err != nil {
		return &[]Client{}, err
	}
	return &clients, err
}

func (c *Client) FindClientByID(db *gorm.DB, uid int32) (*Client, error) {
	// var err error
	err := db.Debug().Model(User{}).Where("id = ?", uid).Take(&c).Error
	if err != nil {
		return &Client{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Client{}, err
	}
	return c, err
}

func (c *Client) UpdateAClient(db *gorm.DB, uid uint32) (*Client, error) {
	db = db.Debug().Model(&Client{}).Where("id = ?", uid).Take(&Client{}).UpdateColumn(
		map[string]interface{}{
			"name":       c.Name,
			"email":      c.Email,
			"updated_at": time.Now(),
		},
	)

	if db.Error != nil {
		return &Client{}, db.Error
	}

	err := db.Debug().Model(&Client{}).Where("id = ?", uid).Take(&c).Error
	if err != nil {
		return &Client{}, err
	}

	return c, nil
}

func (c *Client) DeleteAClient(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&Client{}).Where("id = ?", uid).Take(&Client{}).Delete(&Client{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
