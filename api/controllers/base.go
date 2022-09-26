package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LordCeilan/receiptbackend/api/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=uft8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

		server.DB, err = gorm.Open(Dbdriver, DBURL)

		if err != nil {
			fmt.Printf("cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error: ", err)
		} else {
			fmt.Printf("We are connected to the %s databas", Dbdriver)
		}
	}

	server.DB.Debug().AutoMigrate(&models.User{}, &models.Client{}, &models.Receipt{})
	server.Router = mux.NewRouter()
	// server.initializeRoutes()

}

func (server *Server) Run(addr string) {
	fmt.Printf("Listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
