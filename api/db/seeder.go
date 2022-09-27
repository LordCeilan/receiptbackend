package db

import (
	"log"
	"time"

	"github.com/LordCeilan/receiptbackend/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	{
		ID:       0,
		NickName: "Steven victor",
		Email:    "steven@gmail.com",
		Password: "password",
	},
	{
		NickName: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "password",
	},
}

var clients = []models.Client{
	{
		ID:    0,
		Name:  "Peter Griffin",
		Email: "peter@gmail.com",
	},
	{
		Name:  "Joe Biden",
		Email: "Joe@gmail.com",
	},
}

var receipts = []models.Receipt{
	{
		ID:        0,
		Contents:  "Contents 1",
		Amount:    1000,
		Client:    models.Client{},
		ClientID:  0,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	},
	{
		Contents: "Contents 2",
		Amount:   2000,
	},
}

func Load(db *gorm.DB) {
	var err error
	err = db.Debug().DropTableIfExists(&models.Receipt{}, &models.User{}, &models.Client{}).Error
	if err != nil {
		log.Fatalf("no puede hacer drop a la tabla %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Client{}, &models.Receipt{}).Error
	if err != nil {
		log.Fatalf("no se puede migrar tabla: %v", err)
	}

	err = db.Debug().Model(&models.Receipt{}).AddForeignKey("client_id", "clients(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i := range clients {
		err = db.Debug().Model(&models.Client{}).Create(&clients[i]).Error
		if err != nil {
			log.Fatalf("cannot seed clients table: %v", err)
		}
		receipts[i].ClientID = clients[i].ID

		err = db.Debug().Model(&models.Receipt{}).Create(&receipts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed receipts table: %v", err)
		}
	}

	for j := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[j]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}

	// var err error
	// err = db.Debug().DropTableIfExists(&models.User{}).Error
	// fmt.Println(err)
	// if err != nil {
	// 	log.Fatalf("cannot drop table: %v", err)
	// }

	// for i := range users {
	// 	err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
	// 	if err != nil {
	// 		log.Fatalf("cannot seed users table: %v", err)
	// 	}
	// }

}
