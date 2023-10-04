package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
func GetCard(db *sql.DB, id int) (Card, error) {
	getCardSQL := fmt.Sprintf("SELECT * FROM cards WHERE id='%v';", strconv.Itoa(id))
	rows, err := db.Query(getCardSQL)
	var card Card
	for rows.Next() {
		err = rows.Scan(&card.ID, &card.UID, &card.Comment, &card.Balance, &card.CreationTimestamp, &card.LastTransaction)
		if err != nil {
			return card, err
		}
	}
	return card, nil
}

func GetUserByCookie(r *http.Request) (User, error) {
	c, err := r.Cookie("token")
	if err != nil {
		return User{}, NoCookie
	}
	tknStr := c.Value
	u, ok := CookieToUserMap[tknStr]
	if ok {
		result, err := GetUser(db, u)
		return result, err
	}
	return User{}, NoCookie
}

func CheckCardOwner(r *http.Request) bool {
	vars := mux.Vars(r)
	id := vars["id"]
	u, _ := GetUserByCookie(r)
	c, _ := GetCard(db, id)
	if u.ID == c.UID {
		return true
	} else {
		return false
	}
}
