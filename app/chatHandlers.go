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
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	var chats []Chat
	rows, err := db.Query(fmt.Sprintf("SELECT chats.id, chats.name, chats.description, chats.created_timestamp FROM chats INNER JOIN chat_users ON chats.id = chat_users.cid WHERE chat_users.uid = %v;", user.ID))
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
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	vars := mux.Vars(r)
	id := vars["id"]
	rows, err := db.Query(fmt.Sprintf("SELECT chats.id, chats.name, chats.description, chats.created_timestamp FROM chats INNER JOIN chat_users ON chats.id = chat_users.cid WHERE chats.id = %v AND chat_users.uid = %v;", id, user.ID))
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
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	timestamp := time.Now().Unix()
	var tmp Chat
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := db.Exec(fmt.Sprintf("INSERT INTO chats (name, description, created_timestamp, admin_id) VALUES ('%v', '%v', %v, %v);", tmp.Name, tmp.Description, timestamp, user.ID))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cid, _ := result.LastInsertId()
	_, err = db.Exec(fmt.Sprintf("INSERT INTO chat_users (cid, uid, entry_timestamp) VALUES (%v, %v, %v);", cid, user.ID, timestamp))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func UpdateChat(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	vars := mux.Vars(r)
	id := vars["id"]
	var tmp Chat
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmp.ID, err = strconv.Atoi(id)
	_, err = db.Query(fmt.Sprintf("UPDATE chats SET name='%v',description='%v' WHERE id='%v' AND admin_id = %v;", tmp.Name, tmp.Description, tmp.ID, user.ID))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func DeleteChatByID(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	vars := mux.Vars(r)
	id := vars["id"]
	_, err := db.Exec(fmt.Sprintf("DELETE FROM chats WHERE id = %v AND admin_id=%v;", id, user.ID))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = db.Exec(fmt.Sprintf("DELETE FROM chat_users WHERE cid = %v;", id))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
