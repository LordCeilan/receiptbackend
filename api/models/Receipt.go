package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Receipt struct {
	ClientID  uint64    `gorm:"size:255;not null:unique" json:"clientid"`
	Contents  string    `gorm:"size:255;not null:unique" json:"content"`
	Amount    uint32    `gorm:"not null" json:"amount"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (r *Receipt) Prepare() {
	r.ClientID = 0
	r.Contents = ""
}

func (r *Receipt) SaveReceipt(db *gorm.DB, c *Client) (*Receipt, error) {
	var err error
	err = db.Debug().Model(&Receipt{}).Where("clientid = ?", c.Receipts.ClientID).Create(&c).Error
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Receipt) FindAllReceipts(db *gorm.DB, c *Client) (*[]Receipt, error) {
	var err error
	receipts := []Receipt{}
	err = db.Debug().Model(&Receipt{}).Where("clientid = ?", c.Receipts.ClientID).Find(&receipts).Error

	if err != nil {
		return &[]Receipt{}, nil
	}

	if len(receipts) > 0 {
		for i := range receipts {
			err := db.Debug().Model(&Client{}).Where("receipt = ? ", receipts[i].ClientID).Take(&receipts[i].Contents).Error
			if err != nil {
				return &[]Receipt{}, err
			}
		}
	}
	return &receipts, nil
}

// func (r *Receipt) FindReceiptById(db *gorm.DB, )
