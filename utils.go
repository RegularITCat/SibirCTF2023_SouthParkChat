package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func CountUser(db *sql.DB, username string) (int, error) {
	countUserSQL := fmt.Sprintf("SELECT count(*) FROM users WHERE login='%v';", username)
	rows, err := db.Query(countUserSQL)
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return -1, err
		}
	}
	return count, nil
}

var NoCookie = errors.New("no cookie in request")
var NoUserWithThatCookie = errors.New("no user with that cookie")

func hashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPwd))
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func checkUserCookie(r *http.Request) (string, error) {
	c, err := r.Cookie("token")
	if err != nil {
		return "", NoCookie
	}
	tknStr := c.Value
	u, ok := CookieToUserMap[tknStr]
	if ok {
		return u, nil
	}
	return "", NoUserWithThatCookie
}

func GetUser(db *sql.DB, username string) (User, error) {
	getUserSQL := fmt.Sprintf("SELECT `login`,`password`,`id` FROM `users` WHERE `login`='%v';", username)
	rows, err := db.Query(getUserSQL)
	var user User
	for rows.Next() {
		err = rows.Scan(&user.Login, &user.Password, &user.ID)
		if err != nil {
			return user, err
		}
	}
	return user, nil
}
