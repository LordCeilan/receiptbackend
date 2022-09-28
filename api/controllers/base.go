package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LordCeilan/receiptbackend/api/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql database, stuff never ever works f*@k
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
)

// Server struct manages the database and the mux router to
// work with endtry points and the DB connections
type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

// Initialize method receives each of the databases files to use with the
// instance, receiver is server to ease of the implementation.

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

	if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	}

	// Create the instance received in the router
	// the newRouter uses mux from gorilla mux to manage
	// the endpoints

	server.Router = mux.NewRouter()

	//Server.DB.Debug.Automigrate receives each of the models
	//for the database creation
	//these creates and modifies the database to have at least
	//in postgresql each of the tables well created
	//AutoMigrate run auto migration for given models,
	//will only add missing fields, won't delete/change current data

	server.DB.Debug().AutoMigrate(&models.User{}, &models.Client{}, &models.Receipt{})
	// server.DB.Debug().AutoMigrate(&models.User{}, &models.Client{}, &models.Receipt{})

	//fs uses a fileServer to index the dist files for the static
	//frontend implementation, also this one will work with the front
	//end implementation for serving and queryng
	fs := http.FileServer(http.Dir("../receiptfrontend/dist"))
	server.Router.PathPrefix("/vue").Handler(fs)

	//server route will ease the proccess of receiving and handling the
	//the controlers routes
	http.Handle("/vue", server.Router)
	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Printf("Listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
