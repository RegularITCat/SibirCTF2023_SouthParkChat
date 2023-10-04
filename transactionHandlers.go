package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	var transactions []Transaction
	rows, err := db.Query("SELECT * FROM transactions;")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var transaction Transaction
		err = rows.Scan(&transaction.ID, &transaction.FromCard, &transaction.ToCard, &transaction.Amount, &transaction.Comment, &transaction.Timestamp)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, transaction)
	}
	result, err := json.Marshal(transactions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM transactions WHERE id = %v", id))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var transaction Transaction
	for rows.Next() {
		err = rows.Scan(&transaction.ID, &transaction.FromCard, &transaction.ToCard, &transaction.Amount, &transaction.Comment, &transaction.Timestamp)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
	result, err := json.Marshal(transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var tmp Transaction
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = db.Exec(fmt.Sprintf("INSERT INTO transactions (from_card, to_card, amount, comment, timestamp) VALUES ('%v', '%v', %v, '%v', %v);", tmp.FromCard, tmp.ToCard, tmp.Amount, tmp.Comment, time.Now().Unix()))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
