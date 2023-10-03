package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, err := json.Marshal(map[string]string{"data": "backend is alive."})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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
			w.WriteHeader(http.StatusOK)
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
		w.WriteHeader(http.StatusOK)
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
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}

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

func GetChats(w http.ResponseWriter, r *http.Request) {
	var chats []Chat
	rows, err := db.Query("SELECT * FROM chats;")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var chat Chat
		err = rows.Scan(&chat.ID, &chat.Name, &chat.Description, &chat.CreatedTimestamp)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		chats = append(chats, chat)
	}
	result, err := json.Marshal(chats)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func GetChatByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM chats WHERE id = %v", id))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var chat Chat
	for rows.Next() {
		err = rows.Scan(&chat.ID, &chat.Name, &chat.Description, &chat.CreatedTimestamp)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
	result, err := json.Marshal(chat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func CreateChat(w http.ResponseWriter, r *http.Request) {
	var tmp Chat
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = db.Exec(fmt.Sprintf("INSERT INTO chats (name, description, created_timestamp) VALUES ('%v', '%v', %v);", tmp.Name, tmp.Description, time.Now().Unix()))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func UpdateChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var tmp Chat
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmp.ID, err = strconv.Atoi(id)
	_, err = db.Query(fmt.Sprintf("UPDATE chats SET name='%v',description='%v' WHERE id='%v';", tmp.Name, tmp.Description, tmp.ID))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func DeleteChatByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, err := db.Exec(fmt.Sprintf("DELETE FROM chats WHERE id = '%v';", id))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func GetCards(w http.ResponseWriter, r *http.Request) {
	var cards []Card
	rows, err := db.Query("SELECT ID,UID FROM cards;")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var card Card
		err = rows.Scan(&card.ID, &card.UID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cards = append(cards, card)
	}
	result, err := json.Marshal(cards)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func GetCardByID(w http.ResponseWriter, r *http.Request) {
}

func CreateCard(w http.ResponseWriter, r *http.Request) {
}

func UpdateCard(w http.ResponseWriter, r *http.Request) {
}

func DeleteCard(w http.ResponseWriter, r *http.Request) {
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
}

func GetMessageByID(w http.ResponseWriter, r *http.Request) {
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
}

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
}

func GetFiles(w http.ResponseWriter, r *http.Request) {
}

func GetFileByID(w http.ResponseWriter, r *http.Request) {
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {
}

func GetTransactionByID(w http.ResponseWriter, r *http.Request) {
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
}

func GetPostByID(w http.ResponseWriter, r *http.Request) {
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
}
