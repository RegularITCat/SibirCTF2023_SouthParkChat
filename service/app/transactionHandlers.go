package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.Atoi(vars["cid"])
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	ownerOk, err := CheckCardOwner(r)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	if !ownerOk {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	}
	transactions, err := GetTransactions(cid, cid)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	result, err := json.Marshal(transactions)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func GetTransactionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	cid, err := strconv.Atoi(vars["cid"])
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	//TODO check card owner
	transaction, err := GetTransaction(id, cid)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	result, err := json.Marshal(transaction)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.Atoi(vars["cid"])
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	ownerOk, err := CheckCardOwner(r)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	if !ownerOk {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	}
	var tmp Transaction
	err = json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	tid, err := CreateTransaction(cid, tmp.ToCard, tmp.Amount, tmp.Comment)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	resultJSON, err := json.Marshal(tid)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resultJSON)
}
