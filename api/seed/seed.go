package seed

import (
	"fmt"
	"log"
	"time"

	"github.com/LordCeilan/receiptbackend/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	{
		ID:        0,
		NickName:  "Steven victor",
		Email:     "steven@gmail.com",
		Password:  "password",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	},
	{
		NickName: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "password",
	},
}

func Load(db *gorm.DB) {
	var err error
	err = db.Debug().DropTableIfExists(&models.User{}).Error
	fmt.Println(err)
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}

}
