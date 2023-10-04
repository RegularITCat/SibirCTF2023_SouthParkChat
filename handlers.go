package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	result, err := json.Marshal(map[string]string{"data": "backend is alive."})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var tmp User
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c, err := CountUser(db, tmp.Login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if c == 1 {
		res, err := GetUser(db, tmp.Login)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if res.Password == tmp.Password {
			cookiePass := hashAndSalt(tmp.Password)
			oldCookie, ok := UserToCookieMap[tmp.Login]
			if ok {
				delete(CookieToUserMap, oldCookie)
				delete(UserToCookieMap, tmp.Login)
			}
			CookieToUserMap[cookiePass] = tmp.Login
			UserToCookieMap[tmp.Login] = cookiePass
			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    cookiePass,
				SameSite: http.SameSiteNoneMode,
				Path:     "/",
				Secure:   true,
			})
			_, err = db.Exec(fmt.Sprintf("UPDATE users SET status = 'online' WHERE login = '%v';", tmp.Login))
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var tmp User
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c, err := CountUser(db, tmp.Login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if c == 0 {
		insertUserSQL := fmt.Sprintf("INSERT INTO users (login, password, first_name, second_name, registration_timestamp, login_timestamp, status) VALUES ('%v', '%v', '%v', '%v', '%v', '%v', '%v');", tmp.Login, tmp.Password, tmp.FirstName, tmp.SecondName, time.Now().Unix(), time.Now().Unix(), "online")
		_, err := db.Exec(insertUserSQL)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		cookiePass := hashAndSalt(tmp.Password)
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    cookiePass,
			SameSite: http.SameSiteNoneMode,
			Path:     "/",
			Secure:   true,
		})
		CookieToUserMap[cookiePass] = tmp.Login
		UserToCookieMap[tmp.Login] = cookiePass
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tknStr := c.Value
	u, ok := CookieToUserMap[tknStr]
	if ok {
		delete(CookieToUserMap, tknStr)
		delete(UserToCookieMap, u)
		_, err = db.Exec(fmt.Sprintf("UPDATE users SET status = 'offline' WHERE login = '%v';", u))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}
