package controllers

import (
	"net/http"

	"github.com/LordCeilan/receiptbackend/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To this awsome API")
}
