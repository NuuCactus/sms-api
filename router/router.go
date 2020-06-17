package router

import (
	"github.com/nuucactus/sms-api/endpoints/messages"
	"github.com/nuucactus/sms-api/endpoints/metrics"

	"github.com/gorilla/mux"
)

func NewRouter() (r *mux.Router) {
	r = mux.NewRouter()
	r.HandleFunc("/messages", messages.GetMessages).Methods("GET")
	r.HandleFunc("/messages", messages.PostMessages).Methods("POST")
	r.HandleFunc("/metrics", metrics.GetMetrics).Methods("GET")
	return r
}
