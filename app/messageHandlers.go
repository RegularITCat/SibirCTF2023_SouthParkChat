package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserByCookie(r)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	vars := mux.Vars(r)
	cid, err := strconv.Atoi(vars["cid"])
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	isIn, err := CheckUserInDB(user.ID, cid)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	if !isIn {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	messages := make([]Message, 0)
	rows, err := db.Query(fmt.Sprintf("SELECT id,cid,uid,message,timestamp FROM messages WHERE cid=%v;", cid))
	defer rows.Close()
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	for rows.Next() {
		var message Message
		err = rows.Scan(&message.ID, &message.CID, &message.UID, &message.Message, &message.Timestamp)
		if err != nil {
			printError(w, r, err, http.StatusInternalServerError)
		}
		messages = append(messages, message)
	}
	result, err := json.Marshal(messages)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserByCookie(r)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	vars := mux.Vars(r)
	cid, err := strconv.Atoi(vars["cid"])
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	isIn, err := CheckUserInDB(user.ID, cid)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	if !isIn {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var tmp Message
	err = json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	result, err := db.Exec(fmt.Sprintf("INSERT INTO messages (cid, uid, message, timestamp) VALUES (%v, %v, '%v', %v);", cid, user.ID, tmp.Message, time.Now().Unix()))
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	mid, err := result.LastInsertId()
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	resultJson, err := json.Marshal(mid)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resultJson)
}

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserByCookie(r)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	vars := mux.Vars(r)
	cid, err := strconv.Atoi(vars["cid"])
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	mid, err := strconv.Atoi(vars["mid"])
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	isIn, err := CheckUserInDB(user.ID, cid)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	if !isIn {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var tmp Message
	err = json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	_, err = db.Query(fmt.Sprintf("UPDATE messages SET message='%v' WHERE id=%v AND uid=%v AND cid=%v;", tmp.Message, mid, user.ID, cid))
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserByCookie(r)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	vars := mux.Vars(r)
	cid, err := strconv.Atoi(vars["cid"])
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	mid, err := strconv.Atoi(vars["mid"])
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	isIn, err := CheckUserInDB(user.ID, cid)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	if !isIn {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	_, err = db.Exec(fmt.Sprintf("DELETE FROM messages WHERE id = %v AND cid=%v AND uid=%v;", mid, cid, user.ID))
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
