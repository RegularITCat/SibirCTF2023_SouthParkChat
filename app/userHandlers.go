package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []User
	rows, err := db.Query("SELECT id, first_name, second_name FROM users;")
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.FirstName, &user.SecondName)
		if err != nil {
			printError(w, r, err, http.StatusInternalServerError)
		}
		users = append(users, user)
	}
	result, err := json.Marshal(users)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserByCookie(r)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = %v", user.ID))
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Login, &user.Password, &user.FirstName, &user.SecondName, &user.RegistrationTimestamp, &user.LoginTimestamp, &user.Status)
		if err != nil {
			printError(w, r, err, http.StatusInternalServerError)
		}
	}
	result, err := json.Marshal(user)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserByCookie(r)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	var tmp User
	err = json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	_, err = db.Query(fmt.Sprintf("UPDATE users SET login='%v',first_name='%v',second_name='%v' WHERE id='%v';", tmp.Login, tmp.FirstName, tmp.SecondName, user.ID))
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	tkn := UserToCookieMap[user.Login]
	delete(CookieToUserMap, user.Login)
	delete(UserToCookieMap, tkn)
	CookieToUserMap[tkn] = tmp.Login
	UserToCookieMap[tmp.Login] = tkn
	w.WriteHeader(http.StatusOK)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserByCookie(r)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	err = DeleteMyUserInDB(user.ID)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
