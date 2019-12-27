package ActivityAPI


import (
	"github.com/gorilla/mux"

	"github.com/buyubaya/go-blog/interfaces"
)


func SetActivityRouter(a interfaces.IApp, s *mux.Router) {
	s.HandleFunc("/", GetActivities(a)).Methods("PUT")
	s.HandleFunc("/{activityID}", DeleteActivity(a)).Methods("DELETE")
}