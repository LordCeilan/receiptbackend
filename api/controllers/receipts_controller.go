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

func (server *Server) CreateReceipt(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	receipt := models.Receipt{}
	err = json.Unmarshal(body, &receipt)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	receipt.Prepare()
	err = receipt.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	_, err = auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// if uid != receipt.ClientID {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	// 	return
	// }

	receiptCreated, err := receipt.SaveReceipt(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, receiptCreated.ID))
	responses.JSON(w, http.StatusCreated, receiptCreated)
}

func (server *Server) GetReceipts(w http.ResponseWriter, r *http.Request) {
	receipt := models.Receipt{}
	receipts, err := receipt.FindAllReceipts(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, receipts)
}

func (server *Server) GetReceipt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	receipt := models.Receipt{}

	receiptReceived, err := receipt.FindReceiptById(server.DB, rid)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, receiptReceived)

}

func (server *Server) UpdateReceipt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check if the receipt id is valid
	rid, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//Check if the auth token is valid and get the user id from it

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	receipt := models.Receipt{}
	err = server.DB.Debug().Model(models.Receipt{}).Where("id = ?", rid).Take(&receipt).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("post not found"))
		return
	}

	//If an client attemp to update a a post not belonging to him

	if uid != receipt.ClientID {
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
	receiptUpdate := models.Receipt{}
	err = json.Unmarshal(body, &receiptUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if uid != receiptUpdate.ClientID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	}

	receiptUpdate.Prepare()
	err = receipt.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	receiptUpdate.ID = receipt.ID

	receiptUpdated, err := receiptUpdate.UpdateAReceipt(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, receiptUpdated)

}

func (server *Server) DeleteAClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Is a valid receipt id given to us?
	cid, err := strconv.ParseUint(vars["id"], 10, 64)

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

	err = server.DB.Debug().Model(models.Receipt{}).Where("id = ?", cid).Take(&receipt).Error

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

	_, err = receipt.DeleteAReceipt(server.DB, cid, uid)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", cid))
	responses.JSON(w, http.StatusNoContent, "")

}
