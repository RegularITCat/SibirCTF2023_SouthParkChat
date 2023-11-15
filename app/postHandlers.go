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

func GetPosts(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var err error
	var page, pageSize int
	pageString, ok := r.Form["page"]
	if !ok {
		page = 1
	} else {
		page, err = strconv.Atoi(pageString[0])
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	pageSizeString, ok := r.Form["pageSize"]
	if !ok {
		pageSize = 10
	} else {
		pageSize, err = strconv.Atoi(pageSizeString[0])
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	posts := make([]Post, 0)
	rows, err := db.Query(fmt.Sprintf("SELECT id,uid,name,creation_timestamp FROM posts LIMIT %v OFFSET %v;", page*pageSize, pageSize))
	defer rows.Close()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.ID, &post.UID, &post.Name, &post.CreationTimestamp)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}
	resultJSON, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resultJSON)
}

func GetPostByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var post Post
	rows, err := db.Query(fmt.Sprintf("SELECT id,uid,content,name,creation_timestamp FROM posts WHERE id=%v;", id))
	defer rows.Close()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		err = rows.Scan(&post.ID, &post.UID, &post.Content, &post.Name, &post.CreationTimestamp)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	resultJSON, err := json.Marshal(post)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resultJSON)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	var tmp Post
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, err = db.Query(fmt.Sprintf("UPDATE posts SET name='%v',content='%v' WHERE id=%v AND uid=%v;", tmp.Name, tmp.Content, tmp.ID, user.ID))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, err = db.Exec(fmt.Sprintf("DELETE FROM posts WHERE id=%v AND uid=%v;", id, user.ID))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	username := CookieToUserMap[c.Value]
	user, _ := GetUser(db, username)
	var tmp Post
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, err = db.Exec(fmt.Sprintf("INSERT INTO posts (uid,name,content,creation_timestamp) VALUES (%v, '%v', '%v', %v);", user.ID, tmp.Name, tmp.Content, time.Now().Unix()))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
