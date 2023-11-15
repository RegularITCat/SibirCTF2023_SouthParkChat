package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.Atoi(vars["cid"])
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	if !CheckCardOwner(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	transactions := make([]Transaction, 0)
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM transactions WHERE from_card=%v OR to_card=%v;", cid, cid))
	defer rows.Close()
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	for rows.Next() {
		var transaction Transaction
		err = rows.Scan(&transaction.ID, &transaction.FromCard, &transaction.ToCard, &transaction.Amount, &transaction.Comment, &transaction.Timestamp)
		if err != nil {
			printError(w, r, err, http.StatusInternalServerError)
		}
		transactions = append(transactions, transaction)
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
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM transactions WHERE id=%v AND from_card=%v;", id, cid))
	defer rows.Close()
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	var transaction Transaction
	for rows.Next() {
		err = rows.Scan(&transaction.ID, &transaction.FromCard, &transaction.ToCard, &transaction.Amount, &transaction.Comment, &transaction.Timestamp)
		if err != nil {
			printError(w, r, err, http.StatusInternalServerError)
		}
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
	if !CheckCardOwner(r) {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	var tmp Transaction
	err = json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	var from_card Card
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM cards WHERE id=%v", cid))
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&from_card.ID, &from_card.UID, &from_card.Comment, &from_card.Balance, &from_card.CreationTimestamp, &from_card.LastTransaction)
		if err != nil {
			printError(w, r, err, http.StatusInternalServerError)
		}
	}
	var to_card Card
	rows, err = db.Query(fmt.Sprintf("SELECT * FROM cards WHERE id=%v", tmp.ToCard))
	for rows.Next() {
		err = rows.Scan(&to_card.ID, &to_card.UID, &to_card.Comment, &to_card.Balance, &to_card.CreationTimestamp, &to_card.LastTransaction)
		if err != nil {
			printError(w, r, err, http.StatusInternalServerError)
		}
	}
	if from_card.Balance-tmp.Amount < 0 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	result, err := db.Exec(fmt.Sprintf("INSERT INTO transactions (from_card, to_card, amount, comment, timestamp) VALUES ('%v', '%v', %v, '%v', %v);", cid, tmp.ToCard, tmp.Amount, tmp.Comment, time.Now().Unix()))
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	tid, err := result.LastInsertId()
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	_, err = db.Query(fmt.Sprintf("UPDATE cards SET balance=%v WHERE id=%v;", from_card.Balance-tmp.Amount, from_card.ID))
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	_, err = db.Query(fmt.Sprintf("UPDATE cards SET balance=%v WHERE id=%v;", to_card.Balance+tmp.Amount, to_card.ID))
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
