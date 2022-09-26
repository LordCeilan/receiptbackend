package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LordCeilan/receiptbackend/api/controllers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func Run() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("../receiptfrontend/dist"))

	r.PathPrefix("/").Handler(fs)
	http.Handle("/", r)

	// server.Initialize(os.Getenv(""), os.Getenv(""), os.Getenv(""), os.Getenv(""), os.Getenv(""), os.Getenv(""))
	// server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	server.Run(":3000")
}
