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

func GetMessages(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	vars := mux.Vars(r)
	cid, err := strconv.Atoi(vars["cid"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	isIn, err := CheckUserInDB(user.ID, cid)
	if !isIn {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	messages := make([]Message, 0)
	rows, err := db.Query(fmt.Sprintf("SELECT id,cid,uid,message,timestamp FROM messages WHERE cid=%v;", cid))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var message Message
		err = rows.Scan(&message.ID, &message.CID, &message.UID, &message.Message, &message.Timestamp)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		messages = append(messages, message)
	}
	result, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	vars := mux.Vars(r)
	cid, err := strconv.Atoi(vars["cid"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	isIn, err := CheckUserInDB(user.ID, cid)
	if !isIn {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var tmp Message
	err = json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, err = db.Exec(fmt.Sprintf("INSERT INTO messages (cid, uid, message, timestamp) VALUES (%v, %v, '%v', %v);", cid, user.ID, tmp.Message, time.Now().Unix()))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	vars := mux.Vars(r)
	cid, err := strconv.Atoi(vars["cid"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	mid, err := strconv.Atoi(vars["mid"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	isIn, err := CheckUserInDB(user.ID, cid)
	if !isIn {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var tmp Message
	err = json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, err = db.Query(fmt.Sprintf("UPDATE messages SET message='%v' WHERE id=%v AND uid=%v AND cid=%v;", tmp.Message, mid, user.ID, cid))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	vars := mux.Vars(r)
	cid, err := strconv.Atoi(vars["cid"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	mid, err := strconv.Atoi(vars["mid"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	isIn, err := CheckUserInDB(user.ID, cid)
	if !isIn {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	_, err = db.Exec(fmt.Sprintf("DELETE FROM messages WHERE id = %v AND cid=%v AND uid=%v;", mid, cid, user.ID))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
