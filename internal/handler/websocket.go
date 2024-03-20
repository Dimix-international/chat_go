package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Dimix-international/chat_go/internal/model"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]string)
	broadcast = make(chan model.WSMessage)
	wsUP      = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.URL.Path, "/ws/")
	if _, ok := model.UserData.TKx[token]; !ok {
		log.Printf("User with token %s not found", token)
		return
	}

	conn, err := wsUP.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	clients[conn] = token

	for {
		_, msgByte, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(msgByte)
			delete(clients, conn)
			return
		}

		broadcast <- model.WSMessage{Text: msgByte, Token: token}
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		user := model.UserData.TKx[msg.Token]
		newMsg := model.Message{
			Text:     string(msg.Text),
			UserID:   user.ID,
			UserName: user.Name,
			Created:  time.Now().UTC().Format(time.DateTime),
		}

		msgJson, err := json.Marshal(newMsg)
		if err != nil {
			log.Println(err)
		}
		//write message in file
		if _, err := fmt.Fprintln(model.DBMessage, string(msgJson)); err != nil {
			log.Println(err)
		}

		for client := range clients {
			if msg.Token != clients[client] {
				if err := client.WriteJSON(newMsg); err != nil {
					fmt.Println(err)
					delete(clients, client)
					client.Close()
				}
			}
		}
	}
}
