package main

type User struct {
	ID                    int    `json:"id,omitempty"`
	Login                 string `json:"login,omitempty"`
	Password              string `json:"password,omitempty"`
	FirstName             string `json:"first_name,omitempty"`
	SecondName            string `json:"second_name,omitempty"`
	RegistrationTimestamp int64  `json:"registration_timestamp,omitempty"`
	LoginTimestamp        int64  `json:"login_timestamp,omitempty"`
	Status                string `json:"status,omitempty"`
}

type Chat struct {
	ID               int    `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	CreatedTimestamp int64  `json:"created_timestamp,omitempty"`
}

type ChatUsers struct {
	ID             int    `json:"id,omitempty"`
	CID            string `json:"cid,omitempty"`
	UID            string `json:"uid,omitempty"`
	EntryTimestamp int64  `json:"entry_timestamp,omitempty"`
}

type Message struct {
	ID        int    `json:"id,omitempty"`
	CID       string `json:"cid,omitempty"`
	RID       string `json:"rid,omitempty"`
	Message   string `json:"message,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type Card struct {
	ID                int     `json:"id,omitempty"`
	UID               string  `json:"uid,omitempty"`
	Comment           string  `json:"comment,omitempty"`
	Balance           float64 `json:"balance,omitempty"`
	CreationTimestamp int64   `json:"creation_timestamp,omitempty"`
	LastTransaction   int64   `json:"last_transaction,omitempty"`
}

type Transaction struct {
	ID        int     `json:"id,omitempty"`
	FromCard  string  `json:"from_card,omitempty"`
	ToCard    string  `json:"to_card,omitempty"`
	Amount    float64 `json:"amount,omitempty"`
	Comment   string  `json:"comment,omitempty"`
	Timestamp int64   `json:"timestamp,omitempty"`
}

type File struct {
	ID              int    `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Path            string `json:"path,omitempty"`
	UploadTimestamp int64  `json:"upload_timestamp,omitempty"`
}