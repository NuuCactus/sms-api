package messages

import (
	"net/http"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get /messages\n"))
}

func PostMessages(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("post /messages\n"))
}
