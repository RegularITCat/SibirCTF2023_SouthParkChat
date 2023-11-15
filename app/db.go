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
		//TODO when date is come, stop giving money for free (p.s. after ctf for example lol)
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

func GetCardsByUserID(userID int) ([]Card, error) {
	cards := make([]Card, 0)
	rows, err := db.Query(fmt.Sprintf("SELECT id,uid,comment,balance,creation_timestamp,last_transaction FROM cards WHERE uid=%v;", userID))
	defer rows.Close()
	if err != nil {
		return cards, err
	}
	for rows.Next() {
		var card Card
		err = rows.Scan(&card.ID, &card.UID, &card.Comment, &card.Balance, &card.CreationTimestamp, &card.LastTransaction)
		if err != nil {
			return cards, err
		}
		cards = append(cards, card)
	}
	return cards, err
}

func GetCardByUserIDAndID(userID, ID int) (Card, error) {
	var card Card
	rows, err := db.Query(fmt.Sprintf("SELECT id,uid,comment,balance,creation_timestamp,last_transaction FROM cards WHERE uid=%v AND id=%v;", userID, ID))
	defer rows.Close()
	if err != nil {
		return card, err
	}
	for rows.Next() {
		err = rows.Scan(&card.ID, &card.UID, &card.Comment, &card.Balance, &card.CreationTimestamp, &card.LastTransaction)
		if err != nil {
			return card, err
		}
	}
	return card, err
}

func CreateCard(userID int, comment string) (int, error) {
	result, err := db.Exec(fmt.Sprintf("INSERT INTO cards (uid, comment, balance, creation_timestamp, last_transaction) VALUES (%v, '%v', %v, %v, %v);", userID, comment, 0, time.Now().Unix(), 0))
	if err != nil {
		return 0, err
	}
	cid, err := result.LastInsertId()
	return int(cid), err
}

func UpdateCard(userID, id int, comment string) error {
	_, err := db.Query(fmt.Sprintf("UPDATE cards SET comment='%v' WHERE id=%v AND uid=%v;", comment, id, userID))
	return err
}

func DeleteCard(userID, id int) error {
	_, err := db.Exec(fmt.Sprintf("DELETE FROM cards WHERE id=%v AND uid=%v;", id, userID))
	return err
}

func GetChats(userID int) ([]Chat, error) {
	chats := make([]Chat, 0)
	rows, err := db.Query(fmt.Sprintf("SELECT chats.id, chats.name, chats.description, chats.created_timestamp FROM chats INNER JOIN chat_users ON chats.id = chat_users.cid WHERE chat_users.uid = %v;", userID))
	defer rows.Close()
	if err != nil {
		return chats, err
	}
	for rows.Next() {
		var chat Chat
		err = rows.Scan(&chat.ID, &chat.Name, &chat.Description, &chat.CreatedTimestamp)
		if err != nil {
			return chats, err
		}
		chats = append(chats, chat)
	}
	return chats, err
}

func GetChat(id, userID int) (Chat, error) {
	var chat Chat
	rows, err := db.Query(fmt.Sprintf("SELECT chats.id, chats.name, chats.description, chats.created_timestamp FROM chats INNER JOIN chat_users ON chats.id = chat_users.cid WHERE chats.id = %v AND chat_users.uid = %v;", id, userID))
	if err != nil {
		return chat, err
	}
	for rows.Next() {
		err = rows.Scan(&chat.ID, &chat.Name, &chat.Description, &chat.CreatedTimestamp)
		if err != nil {
			return chat, err
		}
	}
	return chat, nil
}

func CreateChat(userID int, name, description string) (int, error) {
	timestamp := time.Now().Unix()
	result, err := db.Exec(fmt.Sprintf("INSERT INTO chats (name, description, created_timestamp, admin_id) VALUES ('%v', '%v', %v, %v);", name, description, timestamp, userID))
	if err != nil {
		return 0, err
	}
	cid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	_, err = db.Exec(fmt.Sprintf("INSERT INTO chat_users (cid, uid, entry_timestamp) VALUES (%v, %v, %v);", cid, userID, timestamp))
	return int(cid), err
}

func UpdateChat(id, userID int, name, description string) error {
	_, err := db.Query(fmt.Sprintf("UPDATE chats SET name='%v',description='%v' WHERE id='%v' AND admin_id = %v;", name, description, id, userID))
	return err
}

func DeleteChat(id, userID int) error {
	_, err := db.Exec(fmt.Sprintf("DELETE FROM chats WHERE id = %v AND admin_id=%v;", id, userID))
	if err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf("DELETE FROM chat_users WHERE cid = %v;", id))
	return err
}
