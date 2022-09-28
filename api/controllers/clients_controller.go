package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/LordCeilan/receiptbackend/api/auth"
	"github.com/LordCeilan/receiptbackend/api/models"
	"github.com/LordCeilan/receiptbackend/api/responses"
	"github.com/LordCeilan/receiptbackend/api/util/formaterror"
	"github.com/gorilla/mux"
)

func (server *Server) CreateClient(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	client := models.Client{}
	err = json.Unmarshal(body, &client)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	client.Prepare()
	err = client.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	_, err = auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// if uid != client.ID {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	// 	return
	// }

	clientCreated, err := client.SaveClient(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, clientCreated.ID))
	responses.JSON(w, http.StatusCreated, clientCreated)
}

func (server *Server) GetClients(w http.ResponseWriter, r *http.Request) {
	client := models.Client{}
	clients, err := client.FindAllClients(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, clients)
}

func (server *Server) GetClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.ParseUint(vars["id"], 10, 64)

	fmt.Println(vars)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	client := models.Client{}

	receiptReceived, err := client.FindClientByID(server.DB, cid)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, receiptReceived)

}

func (server *Server) UpdateAClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check if the receipt id is valid
	rid, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//Check if the auth token is valid and get the user id from it

	cid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	client := models.Client{}
	err = server.DB.Debug().Model(models.Client{}).Where("id = ?", rid).Take(&client).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("post not found"))
		return
	}

	//If an client attemp to update a a post not belonging to him

	if cid != client.ID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	//Read the data posted
	body, err := io.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// start processing the request data
	clientUpdate := models.Client{}
	err = json.Unmarshal(body, &clientUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if cid != clientUpdate.ID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	}

	clientUpdate.Prepare()
	err = client.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	clientUpdate.ID = client.ID

	receiptUpdated, err := clientUpdate.UpdateAClient(server.DB, cid)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, receiptUpdated)

}

func (server *Server) DeleteReceipt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Is a valid receipt id given to us?
	rid, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	}

	//checks if receipt exist
	receipt := models.Receipt{}

	err = server.DB.Debug().Model(models.Receipt{}).Where("id = ?", rid).Take(&receipt).Error

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	//Is the authenticated user, the owner of this receipt
	// I need to get this variable
	if uid != receipt.ClientID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	_, err = receipt.DeleteAReceipt(server.DB, rid, uid)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", rid))
	responses.JSON(w, http.StatusNoContent, "")

}
