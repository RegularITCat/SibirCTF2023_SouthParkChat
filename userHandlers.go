package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func GetMyUser(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = %v", user.ID))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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

func UpdateMyUser(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	var tmp User
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = db.Query(fmt.Sprintf("UPDATE users SET login='%v',first_name='%v',second_name='%v' WHERE id='%v';", tmp.Login, tmp.FirstName, tmp.SecondName, user.ID))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tkn := UserToCookieMap[user.Login]
	delete(CookieToUserMap, user.Login)
	delete(UserToCookieMap, tkn)
	CookieToUserMap[tkn] = tmp.Login
	UserToCookieMap[tmp.Login] = tkn
	w.WriteHeader(http.StatusOK)
}

func DeleteMyUser(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	err := DeleteMyUserInDB(user.ID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
