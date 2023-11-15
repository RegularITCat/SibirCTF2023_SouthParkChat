package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var tmp User
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	c, err := CountUser(db, tmp.Login)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if c == 1 {
		res, err := GetUser(db, tmp.Login)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
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
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			return
		} else {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	} else {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var tmp User
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	c, err := CountUser(db, tmp.Login)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if c == 0 {
		err = CreateUser(tmp.Login, tmp.Password, tmp.FirstName, tmp.SecondName)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	ok := CheckValid(w, r)
	if ok {
		return
	}
	result, err := json.Marshal(map[string]string{"data": "backend is alive."})
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
