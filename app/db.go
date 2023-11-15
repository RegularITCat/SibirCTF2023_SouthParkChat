package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

func CreateDB(path string) (*sql.DB, error) {
	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		return &sql.DB{}, err
	}
	_, _ = sqlDB.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, login TEXT NOT NULL, password TEXT NOT NULL, first_name TEXT, second_name TEXT, registration_timestamp INTEGER NOT NULL, login_timestamp INTEGER, status TEXT);")
	_, _ = sqlDB.Exec("CREATE TABLE IF NOT EXISTS chats (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, description TEXT, created_timestamp INTEGER NOT NULL, admin_id INTEGER NOT NULL);")
	_, _ = sqlDB.Exec("CREATE TABLE IF NOT EXISTS messages (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, cid INTEGER NOT NULL, uid INTEGER NOT NULL, message TEXT NOT NULL, timestamp INTEGER NOT NULL);")
	_, _ = sqlDB.Exec("CREATE TABLE IF NOT EXISTS cards (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, uid INTEGER NOT NULL, comment TEXT, balance REAL NOT NULL, creation_timestamp INTEGER NOT NULL, last_transaction INTEGER NOT NULL);")
	_, _ = sqlDB.Exec("CREATE TABLE IF NOT EXISTS transactions (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, from_card INTEGER NOT NULL, to_card INTEGER NOT NULL, amount REAL NOT NULL, comment TEXT, timestamp INTEGER NOT NULL);")
	_, _ = sqlDB.Exec("CREATE TABLE IF NOT EXISTS chat_users (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, cid INTEGER NOT NULL, uid INTEGER NOT NULL, entry_timestamp INTEGER NOT NULL);")
	_, _ = sqlDB.Exec("CREATE TABLE IF NOT EXISTS files (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, path TEXT NOT NULL, upload_timestamp INTEGER NOT NULL);")
	_, _ = sqlDB.Exec("CREATE TABLE IF NOT EXISTS posts (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, uid INTEGER NOT NULL, name TEXT NOT NULL, comment TEXT NOT NULL, creation_timestamp INTEGER NOT NULL);")
	rows, err := sqlDB.Query("SELECT count(*) FROM chats WHERE name='general';")
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return &sql.DB{}, err
		}
	}
	timestamp := time.Now().Unix()
	if count == 0 {
		_, _ = sqlDB.Exec(
			fmt.Sprintf("INSERT INTO chats (id, name, description, created_timestamp) VALUES (0, '%v', '%v', '%v');", "general", "general chat", timestamp),
		)
	}
	rows, err = sqlDB.Query("SELECT count(*) FROM users WHERE id=0;")
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return &sql.DB{}, err
		}
	}
	if count == 0 {
		_, _ = sqlDB.Exec(
			fmt.Sprintf("INSERT INTO users (id, login, password, registration_timestamp, status) VALUES (%v, '%v', '%v', %v, 'offline');", 0, "admin", "admin", timestamp),
		)
	}
	rows, err = sqlDB.Query("SELECT count(*) FROM chat_users WHERE cid=0 AND uid=0;")
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return &sql.DB{}, err
		}
	}
	if count == 0 {
		_, _ = sqlDB.Exec(
			fmt.Sprintf("INSERT INTO users (id, login, password, first_name, second_name, registration_timestamp, status) VALUES (%v, '%v', '%v', '%v', '%v', %v, 'offline');", 0, "admin", "admin", "admin", "admin", timestamp),
		)
	}
	sqlDB.Exec("UPDATE users SET status = 'offline';")
	return sqlDB, nil
}

func CreateUser(login, password, firstName, secondName string) error {
	timestamp := time.Now().Unix()
	insertUserSQL := fmt.Sprintf(
		"INSERT INTO users (login, password, first_name, second_name, registration_timestamp, login_timestamp, status) VALUES ('%v', '%v', '%v', '%v', '%v', '%v', '%v');",
		login,
		password,
		firstName,
		secondName,
		timestamp,
		timestamp,
		"online",
	)
	result, err := db.Exec(insertUserSQL)
	if err != nil {
		return err
	}
	uid, _ := result.LastInsertId()
	insertCardSQL := fmt.Sprintf(
		"INSERT INTO cards (uid, comment, balance, creation_timestamp, last_transaction) VALUES (%v, '%v', %v, '%v', %v);",
		uid,
		fmt.Sprintf("user %v default card", login),
		//TODO when date is come, stop giving money for free
		100.0,
		timestamp,
		0,
	)
	_, err = db.Exec(insertCardSQL)
	if err != nil {
		return err
	}
	insertChatUsersSQL := fmt.Sprintf(
		"INSERT INTO chat_users (cid, uid, entry_timestamp) VALUES ('%v', '%v', '%v');",
		1,
		uid,
		timestamp,
	)
	_, err = db.Exec(insertChatUsersSQL)
	return err
}

func CheckUserInDB(userID, chatID int) (bool, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT count(*) FROM chat_users WHERE cid=%v AND uid=%v;", chatID, userID))
	if err != nil {
		log.Println(err)
		return false, err
	}
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Println(err)
			return false, err
		}
	}
	if count != 0 {
		return true, nil
	}
	return false, nil
}

func DeleteMyUserInDB(userID int) error {
	_, err := db.Exec(fmt.Sprintf("DELETE FROM users WHERE id = %v;", userID))
	if err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf("DELETE FROM chats WHERE admin_id = %v;", userID))
	if err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf("DELETE FROM chat_users WHERE uid = %v;", userID))
	if err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf("DELETE FROM cards WHERE uid = %v;", userID))
	if err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf("DELETE FROM messages WHERE uid = %v;", userID))
	return err
}
