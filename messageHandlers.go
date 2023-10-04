package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	var messages []Message
	rows, err := db.Query("SELECT id,cid,uid,message,timestamp FROM messages;")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var message Message
		err = rows.Scan(&message.ID, &message.CID, &message.UID, &message.Message, &message.Timestamp)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cards = append(messages, message)
	}
	result, err := json.Marshal(messages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func GetMessageByID(w http.ResponseWriter, r *http.Request) {
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
}

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
}
