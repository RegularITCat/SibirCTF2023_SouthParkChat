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

func GetChats(w http.ResponseWriter, r *http.Request) {
	var chats []Chat
	rows, err := db.Query("SELECT * FROM chats;")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var chat Chat
		err = rows.Scan(&chat.ID, &chat.Name, &chat.Description, &chat.CreatedTimestamp)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		chats = append(chats, chat)
	}
	result, err := json.Marshal(chats)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func GetChatByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM chats WHERE id = %v", id))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var chat Chat
	for rows.Next() {
		err = rows.Scan(&chat.ID, &chat.Name, &chat.Description, &chat.CreatedTimestamp)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
	result, err := json.Marshal(chat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func CreateChat(w http.ResponseWriter, r *http.Request) {
	var tmp Chat
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = db.Exec(fmt.Sprintf("INSERT INTO chats (name, description, created_timestamp) VALUES ('%v', '%v', %v);", tmp.Name, tmp.Description, time.Now().Unix()))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func UpdateChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var tmp Chat
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmp.ID, err = strconv.Atoi(id)
	_, err = db.Query(fmt.Sprintf("UPDATE chats SET name='%v',description='%v' WHERE id='%v';", tmp.Name, tmp.Description, tmp.ID))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func DeleteChatByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, err := db.Exec(fmt.Sprintf("DELETE FROM chats WHERE id = '%v';", id))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
