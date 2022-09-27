package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Receipt struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Contents  string    `gorm:"size:255;not null:unique" json:"content"`
	Amount    uint32    `gorm:"not null" json:"amount"`
	Client    Client    `json:"client"`
	ClientID  uint32    `gorm:"not null" json:"clientId"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (r *Receipt) Prepare() {
	r.ClientID = 0
	r.Contents = html.EscapeString(strings.TrimSpace(r.Contents))
	r.Amount = 0
	r.Client = Client{}
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
}

func (r *Receipt) Validate() error {
	if r.Contents == "" {
		return errors.New("required Contents")
	}

	if r.Amount == 0 {
		return errors.New("required Amount")
	}

	if r.ClientID < 1 {
		return errors.New("required Client")
	}
	return nil
}

func (r *Receipt) SaveReceipt(db *gorm.DB) (*Receipt, error) {
	var err error
	err = db.Debug().Model(&Receipt{}).Create(&r).Error
	if err != nil {
		return &Receipt{}, err
	}

	if r.ID != 0 {
		err = db.Debug().Model(&Client{}).Where("id = ?", r.ClientID).Take(&r.Client).Error
		if err != nil {
			return &Receipt{}, err
		}
	}
	return r, nil
}

func (r *Receipt) FindAllReceipts(db *gorm.DB) (*[]Receipt, error) {
	var err error
	receipts := []Receipt{}

	err = db.Debug().Model(&Receipt{}).Limit(100).Find(&receipts).Error

	if err != nil {
		return &[]Receipt{}, nil
	}

	if len(receipts) > 0 {
		for i := range receipts {
			err := db.Debug().Model(&Client{}).Where("id = ? ", receipts[i].ClientID).Take(&receipts[i].Client).Error
			if err != nil {
				return &[]Receipt{}, err
			}
		}
	}
	return &receipts, nil
}

func (r *Receipt) FindReceiptById(db *gorm.DB, rid uint) (*Receipt, error) {
	err := db.Debug().Model(&Receipt{}).Where("id = ?", rid).Take(&r).Error

	if err != nil {
		return &Receipt{}, err
	}

	if r.ID != 0 {
		err = db.Debug().Model(&Receipt{}).Where("id = ?", r.ClientID).Take(&r.Client).Error
		if err != nil {
			return &Receipt{}, err
		}
	}
	return r, nil
}

func (r *Receipt) UpdateAReceipt(db *gorm.DB) (*Receipt, error) {
	err := db.Debug().Model(&Receipt{}).Where("id = ?", r.ID).Updates(Receipt{Contents: r.Contents, Amount: r.Amount, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Receipt{}, err
	}

	if r.ID != 0 {
		err = db.Debug().Model(&Receipt{}).Where("id = ?", r.ClientID).Take(&r.Client).Error
		if err != nil {
			return &Receipt{}, err
		}
	}
	return r, nil
}

func (r *Receipt) DeleteAReceipt(db *gorm.DB, rid uint64, cid uint32) (int64, error) {
	db = db.Debug().Model(&Receipt{}).Where("id = ? and clientId= ?", rid, cid).Take(&Receipt{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("receipt not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
