package server

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dimix-international/chat_go/internal/handler"
	"github.com/Dimix-international/chat_go/internal/model"
	"github.com/Dimix-international/chat_go/internal/utils"
)

func Run() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := launchServer(); err != nil {
			log.Println("stop server" + err.Error())
			exit <- syscall.SIGTERM
			close(exit)
		}
	}()

	<-exit

	defer model.DBUser.Close()
	defer model.DBMessage.Close()
}

func launchServer() error {
	initFilesDB()
	readAllUsers()

	http.HandleFunc("/", checkService)
	http.HandleFunc("/api/sign-in", handler.SignIn)
	http.HandleFunc("/ws/", handler.HandleConnections)

	go handler.HandleMessages()

	if err := http.ListenAndServe(model.Port, nil); err != nil {
		return err
	}

	return nil
}

func checkService(w http.ResponseWriter, r *http.Request) {
	utils.ResponseString(w, `{"success": true}`)
}

func initFilesDB() {
	var err error

	model.DBUser, err = os.OpenFile("./internal/data/users.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Panic("Error open users file" + err.Error())
		return
	}

	model.DBMessage, err = os.OpenFile("./internal/data/messages.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Panic("Error open users file" + err.Error())
		return
	}
}

func readAllUsers() {
	scanner := bufio.NewScanner(model.DBUser)
	for scanner.Scan() {
		var user model.User
		if err := json.Unmarshal([]byte(scanner.Text()), &user); err == nil {
			model.UserData.Items = append(model.UserData.Items, user)
			model.UserData.IDx[user.ID] = &model.UserData.Items[len(model.UserData.Items)-1]
			model.UserData.TKx[user.Token] = &model.UserData.Items[len(model.UserData.Items)-1]
		}
	}
}
