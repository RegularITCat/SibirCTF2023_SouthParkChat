package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetCardsHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserByCookie(r)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	cards, err := GetCardsByUserID(user.ID)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	result, err := json.Marshal(cards)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func GetCardHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserByCookie(r)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	vars := mux.Vars(r)
	//log.Println(vars["id"])
	cid, err := strconv.Atoi(vars["cid"])
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	ownerOk, err := CheckCardOwner(r)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	if !ownerOk {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	}
	card, err := GetCardByUserIDAndID(user.ID, cid)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	result, err := json.Marshal(card)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func CreateCardHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserByCookie(r)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	var tmp Card
	err = json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	cid, err := CreateCard(user.ID, tmp.Comment)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	resultJson, err := json.Marshal(cid)
	w.WriteHeader(http.StatusOK)
	w.Write(resultJson)
}

func UpdateCardHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserByCookie(r)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	vars := mux.Vars(r)
	cid := vars["cid"]
	ownerOk, err := CheckCardOwner(r)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	if !ownerOk {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	}
	var tmp Card
	err = json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	tmp.ID, err = strconv.Atoi(cid)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	err = UpdateCard(user.ID, tmp.ID, tmp.Comment)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteCardHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserByCookie(r)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	vars := mux.Vars(r)
	cid, err := strconv.Atoi(vars["cid"])
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	err = DeleteCard(cid, user.ID)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
