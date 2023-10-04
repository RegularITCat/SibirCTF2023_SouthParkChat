package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

/*
post /api/v1/login login in telecaht
post /api/v1/register register in telechat
post /api/v1/logout logout in telechat


get /api/v1/user info about my current user
get /api/v1/user/[id] get info about other user
put /api/v1/user fix my user info

get /api/v1/health check if telechat is ok

get /api/v1/chat get all available chats
get /api/v1/chat/[id] get chat by id
post /api/v1/chat create a chat
put /api/v1/chat/[id] update a chat info
delete /api/v1/chat/[id] delete a chat

get /api/v1/card get all cards
get /api/v1/card/[card_id] get info about my card
post /api/v1/card register new card
put /api/v1/card/[card_id] update card info
delete /api/v1/card/[card_id] delete card

get /api/v1/[chat_id]/msg get all messages
get /api/v1/[chat_id]/msg/[id] get message by id
post /api/v1/[chat_id]/msg send message
put /api/v1/[chat_id]/msg/[id] update message
delete /api/v1/[chat_id]/msg/[id] delete message

get /api/v1/file get all files id
get /api/v1/file/[id] get file
post /api/v1/file send file
delete /api/v1/file/[id] delete file

get /api/v1/transaction get your transactions info
get /api/v1/transaction/[id] get one transaction
post /api/v1/transaction create transaction

get /api/v1/[user_id]/post get all users posts
get /api/v1/[user_id]/post/[id] get one exact post
post /api/v1/[user_id]/post create post
put /api/v1/[user_id]/post/[id] update post
delete /api/v1/[user_id]/post/[id] delete post

*/

var db *sql.DB
var CookieToUserMap = make(map[string]string)
var UserToCookieMap = make(map[string]string)

func main() {
	db, _ = CreateDB(":memory:")
	log.Println(db)
	router := mux.NewRouter()
	router.Use(PanicRecoveryMiddleware, LoggingMiddleware, AuthCheckMiddleware)
	router.HandleFunc("/api/v1/health", HealthHandler)
	router.HandleFunc("/api/v1/login", LoginHandler).Methods("POST")
	router.HandleFunc("/api/v1/register", RegisterHandler).Methods("POST")
	router.HandleFunc("/api/v1/logout", LogoutHandler).Methods("POST")
	router.HandleFunc("/api/v1/user/{id:[0-9]+}", GetUserByID).Methods("GET")
	router.HandleFunc("/api/v1/user", GetUsers).Methods("GET")
	router.HandleFunc("/api/v1/user/{id:[0-9]+}", UpdateUser).Methods("PUT")
	router.HandleFunc("/api/v1/chat", GetChats).Methods("GET")
	router.HandleFunc("/api/v1/chat/{id:[0-9]+}", GetChatByID).Methods("GET")
	router.HandleFunc("/api/v1/chat/{id:[0-9]+}", UpdateChat).Methods("PUT")
	router.HandleFunc("/api/v1/chat/{id:[0-9]+}", DeleteChatByID).Methods("DELETE")
	router.HandleFunc("/api/v1/chat", CreateChat).Methods("POST")
	router.HandleFunc("/api/v1/message", GetMessages).Methods("GET")
	router.HandleFunc("/api/v1/message/{id:[0-9]+}", GetMessageByID).Methods("GET")
	router.HandleFunc("/api/v1/message/{id:[0-9]+}", UpdateMessage).Methods("PUT")
	router.HandleFunc("/api/v1/message/{id:[0-9]+}", DeleteMessage).Methods("DELETE")
	router.HandleFunc("/api/v1/message", CreateMessage).Methods("POST")
	router.HandleFunc("/api/v1/card", GetCards).Methods("GET")
	router.HandleFunc("/api/v1/card/{id:[0-9]+}", GetCardByID).Methods("GET")
	router.HandleFunc("/api/v1/card/{id:[0-9]+}", UpdateCard).Methods("PUT")
	router.HandleFunc("/api/v1/card/{id:[0-9]+}", DeleteCardByID).Methods("DELETE")
	router.HandleFunc("/api/v1/card", CreateCard).Methods("POST")
	router.HandleFunc("/api/v1/transaction", GetTransactions).Methods("GET")
	router.HandleFunc("/api/v1/transaction/{id:[0-9]+}", GetTransactionByID).Methods("GET")
	router.HandleFunc("/api/v1/transaction", CreateTransaction).Methods("POST")
	router.HandleFunc("/api/v1/post", GetPosts).Methods("GET")
	router.HandleFunc("/api/v1/post/{id:[0-9]+}", GetPostByID).Methods("GET")
	router.HandleFunc("/api/v1/post/{id:[0-9]+}", UpdatePost).Methods("PUT")
	router.HandleFunc("/api/v1/post/{id:[0-9]+}", DeletePost).Methods("DELETE")
	router.HandleFunc("/api/v1/post", CreatePost).Methods("POST")
	http.Handle("/", router)
	server := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8888",
	}

	log.Println(server.ListenAndServe())
}
