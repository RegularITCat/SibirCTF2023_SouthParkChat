package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
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
	messages, err := GetMessages(cid)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	result, err := json.Marshal(messages)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
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
	result, err := CreateMessage(cid, user.ID, tmp.Message)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	resultJson, err := json.Marshal(result)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resultJson)
}

func UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
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
	err = UpdateMessage(mid, user.ID, cid, tmp.Message)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
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
	err = DeleteMessage(mid, user.ID, cid)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
