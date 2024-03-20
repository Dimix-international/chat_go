package handler

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Dimix-international/chat_go/internal/model"
	"github.com/Dimix-international/chat_go/internal/utils"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		utils.ResponseString(w, `{"success": false, "msg": "POST method is required"}`)
		return
	}

	var user model.User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.ResponseString(w, fmt.Sprintf(`{"success": false, "msg": "%s"}`, err.Error()))
		return
	}
	if err := json.Unmarshal(body, &user); err != nil {
		utils.ResponseString(w, fmt.Sprintf(`{"success": false, "msg": "%s"}`, err.Error()))
		return
	}

	user.Name = strings.Trim(user.Name, " ")

	if user.Name == "" {
		utils.ResponseString(w, `{"success": false, "msg": "Please enter your name"}`)
		return
	}

	hasher := sha1.New()
	hasher.Write([]byte(fmt.Sprintf(`%s_%v`, user.Name, time.Now().UTC())))

	user.Token = hex.EncodeToString(hasher.Sum(nil))
	user.ID = len(model.UserData.Items) + 1

	userJson, _ := json.Marshal(user)
	//Write to end file
	if _, err := fmt.Fprintln(model.DBUser, string(userJson)); err != nil {
		utils.ResponseString(w, fmt.Sprintf(`{"success": false, "msg": "%s"}`, err.Error()))
		return
	}

	model.UserData.Items = append(model.UserData.Items, user)
	model.UserData.IDx[user.ID] = &model.UserData.Items[len(model.UserData.Items)-1]
	model.UserData.TKx[user.Token] = &model.UserData.Items[len(model.UserData.Items)-1]

	utils.ResponseJson(w, userJson)
}
