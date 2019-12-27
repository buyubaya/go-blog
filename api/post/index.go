package PostAPI


import (
	"github.com/gorilla/mux"

	"github.com/buyubaya/go-blog/interfaces"
)


func SetPostRouter(a interfaces.IApp, s *mux.Router) {
	s.HandleFunc("/", GetList(a)).Methods("PUT")
	s.HandleFunc("/{postID}", GetOne(a)).Methods("GET")
	s.HandleFunc("/", CreatePost(a)).Methods("POST")
	s.HandleFunc("/{postID}", UpdatePost(a)).Methods("PUT")
	s.HandleFunc("/{postID}", DeletePost(a)).Methods("DELETE")
	
}