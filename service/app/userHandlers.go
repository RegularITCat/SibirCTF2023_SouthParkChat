package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	page, pageSize, err := parseFormPageParams(r)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	users, err := GetUsers(page, pageSize)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
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
