package main

import (
	"log"
	"net/http"
)

const (
	port string = ":8989"
)

func main() {
	http.HandleFunc("/", checkService)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Panic("error start serve" + err.Error())
	}
}

func checkService(w http.ResponseWriter, r *http.Request) {
	responseString(w, `{"success": true}`)
}

func responseString(w http.ResponseWriter, text string) {
	responseJson(w, []byte(text))
}

func responseJson(w http.ResponseWriter, v []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(v); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Error`))
	}
}
