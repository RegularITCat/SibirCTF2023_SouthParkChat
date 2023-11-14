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
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	var cards []Card
	rows, err := db.Query(fmt.Sprintf("SELECT id,uid,comment,balance,creation_timestamp,last_transaction FROM cards WHERE uid=%v;", user.ID))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var card Card
		err = rows.Scan(&card.ID, &card.UID, &card.Comment, &card.Balance, &card.CreationTimestamp, &card.LastTransaction)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		cards = append(cards, card)
	}
	result, err := json.Marshal(cards)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func GetCardByID(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	if !CheckCardOwner(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	rows, err := db.Query(fmt.Sprintf("SELECT id,uid,comment,balance,creation_timestamp,last_transaction FROM cards WHERE uid=%v AND id=%v;", user.ID, id))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var card Card
	for rows.Next() {
		err = rows.Scan(&card.ID, &card.UID, &card.Comment, &card.Balance, &card.CreationTimestamp, &card.LastTransaction)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	}
	result, err := json.Marshal(card)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func CreateCard(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	var tmp Card
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, err = db.Exec(fmt.Sprintf("INSERT INTO cards (uid, comment, balance, creation_timestamp, last_transaction) VALUES (%v, '%v', %v, %v, %v);", user.ID, tmp.Comment, 0, time.Now().Unix(), 0))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateCard(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	vars := mux.Vars(r)
	id := vars["id"]
	if !CheckCardOwner(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	var tmp Card
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	tmp.ID, err = strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, err = db.Query(fmt.Sprintf("UPDATE cards SET comment='%v' WHERE id=%v AND uid=%v;", tmp.Comment, tmp.ID, user.ID))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteCardByID(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	_, err := db.Exec(fmt.Sprintf("DELETE FROM cards WHERE id=%v AND uid=%v;", id, user.ID))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
