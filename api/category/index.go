package CategoryAPI


import (
	"github.com/gorilla/mux"

	"github.com/buyubaya/go-blog/interfaces"
)


func SetCategoryRouter(a interfaces.IApp, s *mux.Router) {
	s.HandleFunc("/", GetCategories(a)).Methods("GET")
	s.HandleFunc("/", CreateCategory(a)).Methods("POST")
	s.HandleFunc("/{categoryID}", UpdateCategory(a)).Methods("PUT")
	s.HandleFunc("/{categoryID}", DeleteCategory(a)).Methods("DELETE")
}