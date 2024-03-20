package utils

import "net/http"

func ResponseString(w http.ResponseWriter, text string) {
	ResponseJson(w, []byte(text))
}

func ResponseJson(w http.ResponseWriter, v []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(v); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Error`))
	}
}
