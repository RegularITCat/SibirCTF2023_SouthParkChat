package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func CreateDB(path string) (*sql.DB, error) {
	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		return &sql.DB{}, err
	}
	_, _ = sqlDB.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, login TEXT NOT NULL, password TEXT NOT NULL, first_name TEXT, second_name TEXT, registration_timestamp INTEGER NOT NULL, login_timestamp INTEGER, status TEXT);")
	_, _ = sqlDB.Exec("CREATE TABLE IF NOT EXISTS chats (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, description TEXT, created_timestamp INTEGER NOT NULL);")
	_, _ = sqlDB.Exec("CREATE TABLE IF NOT EXISTS messages (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, cid TEXT NOT NULL, uid TEXT NOT NULL, message TEXT NOT NULL, timestamp INTEGER NOT NULL);")
	_, _ = sqlDB.Exec("CREATE TABLE IF NOT EXISTS cards (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, uid TEXT NOT NULL, comment TEXT, balance REAL NOT NULL, creation_timestamp INTEGER NOT NULL, last_transaction INTEGER);")
	_, _ = sqlDB.Exec("CREATE TABLE IF NOT EXISTS transactions (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, from_card TEXT NOT NULL, to_card TEXT NOT NULL, amount REAL NOT NULL, comment TEXT, timestamp INTEGER NOT NULL);")
	_, _ = sqlDB.Exec("CREATE TABLE IF NOT EXISTS chat_users (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, cid TEXT NOT NULL, uid TEXT NOT NULL, entry_timestamp INTEGER NOT NULL);")
	_, _ = sqlDB.Exec("CREATE TABLE IF NOT EXISTS files (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, path TEXT NOT NULL, upload_timestamp INTEGER NOT NULL);")
	return sqlDB, nil
}
