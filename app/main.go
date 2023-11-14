package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var db *sql.DB
var CookieToUserMap = make(map[string]string)
var UserToCookieMap = make(map[string]string)

func main() {
	log.Println("Starting chat service...")
	log.Println("Initializing database...")
	db, _ = CreateDB(":memory:")
	dbAddr, exists := os.LookupEnv("SOUTHPARKCHAT_DB_ADDR")
	if exists {
		db, _ = CreateDB(dbAddr)
	}
	log.Println("...done")
	log.Println("Initializing router...")
	router := mux.NewRouter()
	router.Use(PanicRecoveryMiddleware, LoggingMiddleware, AuthCheckMiddleware)
	router.HandleFunc("/api/v1/health", HealthHandler)
	router.HandleFunc("/api/v1/login", LoginHandler).Methods("POST")
	router.HandleFunc("/api/v1/register", RegisterHandler).Methods("POST")
	router.HandleFunc("/api/v1/logout", LogoutHandler).Methods("POST")
	router.HandleFunc("/api/v1/user/me", GetMyUser).Methods("GET")
	router.HandleFunc("/api/v1/user", GetUsers).Methods("GET")
	router.HandleFunc("/api/v1/user/me", UpdateMyUser).Methods("PUT")
	router.HandleFunc("/api/v1/user/me", DeleteMyUser).Methods("DELETE")
	router.HandleFunc("/api/v1/chat", GetChats).Methods("GET")
	router.HandleFunc("/api/v1/chat/{id:[0-9]+}", GetChatByID).Methods("GET")
	router.HandleFunc("/api/v1/chat/{id:[0-9]+}", UpdateChat).Methods("PUT")
	router.HandleFunc("/api/v1/chat/{id:[0-9]+}", DeleteChatByID).Methods("DELETE")
	router.HandleFunc("/api/v1/chat", CreateChat).Methods("POST")
	router.HandleFunc("/api/v1/chat/{cid:[0-9]+}/message", GetMessages).Methods("GET")
	router.HandleFunc("/api/v1/chat/{cid:[0-9]+}/message/{mid:[0-9]+}", UpdateMessage).Methods("PUT")
	router.HandleFunc("/api/v1/chat/{cid:[0-9]+}/message/{mid:[0-9]+}", DeleteMessage).Methods("DELETE")
	router.HandleFunc("/api/v1/chat/{cid:[0-9]+}/message", CreateMessage).Methods("POST")
	router.HandleFunc("/api/v1/card", GetCards).Methods("GET")
	router.HandleFunc("/api/v1/card/{id:[0-9]+}", GetCardByID).Methods("GET")
	router.HandleFunc("/api/v1/card/{id:[0-9]+}", UpdateCard).Methods("PUT")
	router.HandleFunc("/api/v1/card/{id:[0-9]+}", DeleteCardByID).Methods("DELETE")
	router.HandleFunc("/api/v1/card", CreateCard).Methods("POST")
	router.HandleFunc("/api/v1/card/{cid:[0-9]+}/transaction", GetTransactions).Methods("GET")
	router.HandleFunc("/api/v1/card/{cid:[0-9]+}/transaction/{id:[0-9]+}", GetTransactionByID).Methods("GET")
	router.HandleFunc("/api/v1/card/{cid:[0-9]+}/transaction", CreateTransaction).Methods("POST")
	//router.HandleFunc("/api/v1/post", GetPosts).Methods("GET")
	//router.HandleFunc("/api/v1/post/{id:[0-9]+}", GetPostByID).Methods("GET")
	//router.HandleFunc("/api/v1/post/{id:[0-9]+}", UpdatePost).Methods("PUT")
	//router.HandleFunc("/api/v1/post/{id:[0-9]+}", DeletePost).Methods("DELETE")
	//router.HandleFunc("/api/v1/post", CreatePost).Methods("POST")
	http.Handle("/", router)
	log.Println("...done")
	log.Println("Starting server!")
	server := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8888",
	}
	addr, exists := os.LookupEnv("SOUTHPARKCHAT_ADDR")
	if exists {
		server = &http.Server{
			Handler: router,
			Addr:    addr,
		}
	}

	log.Println(server.ListenAndServe())
}
