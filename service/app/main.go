package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

var db *sql.DB
var CookieToUserMap = make(map[string]string)
var UserToCookieMap = make(map[string]string)

func main() {
	log.Println("Starting chat service...")
	dbAddr, exists := os.LookupEnv("SOUTHPARKCHAT_DB_ADDR")
	if !exists {
		dbAddr = ":memory:"
	} else {
		_, err := os.Stat(dbAddr)
		if os.IsNotExist(err) {
			errDir := os.MkdirAll(filepath.Dir(dbAddr), 0755)
			if errDir != nil {
				log.Fatal(err)
			}
		}
	}
	log.Printf("Initializing database in %v...", dbAddr)
	db, _ = CreateDB(dbAddr)
	log.Println("...done")
	log.Println("Initializing router...")
	router := mux.NewRouter()
	router.Use(PanicRecoveryMiddleware, LoggingMiddleware, AuthCheckMiddleware)
	router.HandleFunc("/api/v1/health", HealthHandler)
	router.HandleFunc("/api/v1/login", LoginHandler).Methods("POST")
	router.HandleFunc("/api/v1/register", RegisterHandler).Methods("POST")
	router.HandleFunc("/api/v1/logout", LogoutHandler).Methods("POST")
	router.HandleFunc("/api/v1/user/me", GetUserHandler).Methods("GET")
	router.HandleFunc("/api/v1/user", GetUsersHandler).Methods("GET")
	router.HandleFunc("/api/v1/user/me", UpdateUserHandler).Methods("PUT")
	router.HandleFunc("/api/v1/user/me", DeleteUserHandler).Methods("DELETE")
	router.HandleFunc("/api/v1/chat", GetChatsHandler).Methods("GET")
	router.HandleFunc("/api/v1/chat/{id:[0-9]+}", GetChatHandler).Methods("GET")
	router.HandleFunc("/api/v1/chat/{cid:[0-9]+}/invite/{uid:[0-9]+}", InviteHandler).Methods("POST")
	router.HandleFunc("/api/v1/chat/{id:[0-9]+}", UpdateChatHandler).Methods("PUT")
	router.HandleFunc("/api/v1/chat/{id:[0-9]+}", DeleteChatHandler).Methods("DELETE")
	router.HandleFunc("/api/v1/chat", CreateChatHandler).Methods("POST")
	router.HandleFunc("/api/v1/chat/{cid:[0-9]+}/message/{mid:[0-9]+}", GetMessageHandler).Methods("GET")
	router.HandleFunc("/api/v1/chat/{cid:[0-9]+}/message", GetMessagesHandler).Methods("GET")
	router.HandleFunc("/api/v1/chat/{cid:[0-9]+}/message/{mid:[0-9]+}", UpdateMessageHandler).Methods("PUT")
	router.HandleFunc("/api/v1/chat/{cid:[0-9]+}/message/{mid:[0-9]+}", DeleteMessageHandler).Methods("DELETE")
	router.HandleFunc("/api/v1/chat/{cid:[0-9]+}/message", CreateMessageHandler).Methods("POST")
	router.HandleFunc("/api/v1/card", GetCardsHandler).Methods("GET")
	router.HandleFunc("/api/v1/card/{cid:[0-9]+}", GetCardHandler).Methods("GET")
	router.HandleFunc("/api/v1/card/{cid:[0-9]+}", UpdateCardHandler).Methods("PUT")
	router.HandleFunc("/api/v1/card/{cid:[0-9]+}", DeleteCardHandler).Methods("DELETE")
	router.HandleFunc("/api/v1/card", CreateCardHandler).Methods("POST")
	router.HandleFunc("/api/v1/card/{cid:[0-9]+}/transaction", GetTransactionsHandler).Methods("GET")
	router.HandleFunc("/api/v1/card/{cid:[0-9]+}/transaction/{id:[0-9]+}", GetTransactionHandler).Methods("GET")
	router.HandleFunc("/api/v1/card/{cid:[0-9]+}/transaction", CreateTransactionHandler).Methods("POST")
	router.HandleFunc("/api/v1/file", GetFilesHandler).Methods("GET")
	router.HandleFunc("/api/v1/file/{id:[0-9]+}", GetFileHandler).Methods("GET")
	router.HandleFunc("/api/v1/file", UploadFileHandler).Methods("POST")
	router.HandleFunc("/api/v1/file/{id:[0-9]+}/download", DownloadFileHandler).Methods("GET")
	router.HandleFunc("/api/v1/file/{id:[0-9]+}", DeleteFileHandler).Methods("DELETE")
	router.HandleFunc("/api/v1/post", GetPostsHandler).Methods("GET")
	router.HandleFunc("/api/v1/post/{id:[0-9]+}", GetPostHandler).Methods("GET")
	router.HandleFunc("/api/v1/post/{id:[0-9]+}", UpdatePostHandler).Methods("PUT")
	router.HandleFunc("/api/v1/post/{id:[0-9]+}", DeletePostHandler).Methods("DELETE")
	router.HandleFunc("/api/v1/post", CreatePostHandler).Methods("POST")
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
