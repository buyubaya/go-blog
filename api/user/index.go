package UserAPI


import (
	"github.com/gorilla/mux"

	"github.com/buyubaya/go-blog/interfaces"
)


func SetUserRouter(a interfaces.IApp, s *mux.Router) {
	s.HandleFunc("/", CreateUser(a)).Methods("POST")
	s.HandleFunc("/", GetUsers(a)).Methods("PUT")
	s.HandleFunc("/{UID}", UpdateUser(a)).Methods("PUT")
	// s.HandleFunc("/{UID}", DeletePost(a)).Methods("DELETE")
}