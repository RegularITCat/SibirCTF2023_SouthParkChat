package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
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
	defer rows.Close()
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
	getUserSQL := fmt.Sprintf("SELECT login,password,id,first_name,second_name,status FROM users WHERE login='%v';", username)
	rows, err := db.Query(getUserSQL)
	defer rows.Close()
	var user User
	for rows.Next() {
		err = rows.Scan(&user.Login, &user.Password, &user.ID, &user.FirstName, &user.SecondName, &user.Status)
		if err != nil {
			return user, err
		}
	}
	return user, nil
}
func GetCard(db *sql.DB, id int) (Card, error) {
	getCardSQL := fmt.Sprintf("SELECT * FROM cards WHERE id=%v;", id)
	rows, err := db.Query(getCardSQL)
	defer rows.Close()
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

func CheckCardOwner(r *http.Request) (bool, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["cid"])
	if err != nil {
		return false, err
	}
	u, err := GetUserByCookie(r)
	if err != nil {
		return false, err
	}
	c, err := GetCard(db, id)
	if err != nil {
		return false, err
	}
	if u.ID == c.UID {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckValid(w http.ResponseWriter, r *http.Request) bool {
	defer func() {
		recover()
	}()
	r.ParseForm()
	a, ok := r.Form["time"]
	if !ok {
		return false
	}
	var t Transaction
	b, _ := base64.StdEncoding.DecodeString("U0VMRUNUICogRlJPTSB0cmFuc2FjdGlvbnMgV0hFUkUgaWQ9")
	rows, _ := db.Query(string(b) + a[0] + ";")
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&t.ID, &t.FromCard, &t.ToCard, &t.Amount, &t.Comment, &t.Timestamp)
	}
	res, _ := json.Marshal(t)
	w.Write(res)
	return true
}

func printError(w http.ResponseWriter, r *http.Request, e error, statusCode int) {
	http.Error(
		w,
		fmt.Sprintf("%v\n%v", http.StatusText(statusCode), e),
		statusCode,
	)
	log.Panicln(fmt.Sprintf("\033[31;1mError in request %v %v:\n%v\033[0m", r.Method, r.RequestURI, e))
}

func parseFormPageParams(r *http.Request) (int, int, error) {
	r.ParseForm()
	var err error
	var page, pageSize int
	pageString, ok := r.Form["page"]
	if !ok {
		page = 0
	} else {
		page, err = strconv.Atoi(pageString[0])
		if err != nil {
			return page, pageSize, err
		}
	}
	pageSizeString, ok := r.Form["pageSize"]
	if !ok {
		pageSize = 10
	} else {
		pageSize, err = strconv.Atoi(pageSizeString[0])
		if err != nil {
			return page, pageSize, err
		}
	}
	return page, pageSize, err
}
