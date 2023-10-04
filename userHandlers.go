package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	rows, err := db.Query("SELECT id, first_name, second_name FROM users;")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.FirstName, &user.SecondName)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	result, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = %v", id))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user User
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Login, &user.Password, &user.FirstName, &user.SecondName, &user.RegistrationTimestamp, &user.LoginTimestamp, &user.Status)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
	result, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var tmp User
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmp.ID, err = strconv.Atoi(id)
	_, err = db.Query(fmt.Sprintf("UPDATE users SET login='%v',first_name='%v',second_name='%v' WHERE id='%v';", tmp.Login, tmp.FirstName, tmp.SecondName))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
