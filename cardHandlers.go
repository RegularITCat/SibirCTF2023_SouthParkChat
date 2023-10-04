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

func GetCards(w http.ResponseWriter, r *http.Request) {
	var cards []Card
	rows, err := db.Query("SELECT ID,UID FROM cards;")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var card Card
		err = rows.Scan(&card.ID, &card.UID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cards = append(cards, card)
	}
	result, err := json.Marshal(cards)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func GetCardByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM card WHERE id = %v", id))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var card Card
	for rows.Next() {
		err = rows.Scan(&card.ID, &card.UID, &card.Comment, &card.Balance, &card.CreationTimestamp, &card.LastTransaction)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
	result, err := json.Marshal(card)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func CreateCard(w http.ResponseWriter, r *http.Request) {
	var tmp Card
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = db.Exec(fmt.Sprintf("INSERT INTO cards (uid, comment, balance, creation_timestamp) VALUES ('%v', '%v', %v, %v);", tmp.UID, tmp.Comment, 0, time.Now().Unix()))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func UpdateCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var tmp Card
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmp.ID, err = strconv.Atoi(id)
	_, err = db.Query(fmt.Sprintf("UPDATE cards SET comment='%v' WHERE id='%v';", tmp.Comment, tmp.ID))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func DeleteCardByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, err := db.Exec(fmt.Sprintf("DELETE FROM cards WHERE id = '%v';", id))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
