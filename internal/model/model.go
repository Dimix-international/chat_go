package model

import "os"

var (
	dbUser    *os.File
	dbMessage *os.File
	userData  = Users{
		IDx:   make(map[int]*User, 0),
		TKx:   make(map[string]*User, 0),
		Items: make([]User, 0, 100),
	}
)

type Message struct {
	Created  string `json:"created"`
	Text     string `json:"text"`
	UserID   int    `json:"user_id"`
	UserName string `json:"name"`
}

type wsMessage struct {
	Text  []byte
	Token string
}

type User struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
	Name  string `json:"name"`
}

type Users struct {
	IDx   map[int]*User
	TKx   map[string]*User
	Items []User
}
