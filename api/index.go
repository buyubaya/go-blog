package API


import (
	"github.com/gorilla/mux"

	"github.com/buyubaya/go-blog/interfaces"
	"github.com/buyubaya/go-blog/middlewares"

	UserAPI "github.com/buyubaya/go-blog/api/user"
	PostAPI "github.com/buyubaya/go-blog/api/post"
	CategoryAPI "github.com/buyubaya/go-blog/api/category"
	ActivityAPI "github.com/buyubaya/go-blog/api/activity"
)


func SetAPIRouter(a interfaces.IApp, s *mux.Router) {

	s.Use(middlewares.FirebaseAuthentication(a))

	// USERS
	userSubrouter := s.PathPrefix("/users").Subrouter().StrictSlash(true)
	UserAPI.SetUserRouter(a, userSubrouter)

	// POSTS
	postSubrouter := s.PathPrefix("/posts").Subrouter().StrictSlash(true)
	PostAPI.SetPostRouter(a, postSubrouter)

	// CATEGORIES
	categorySubrouter := s.PathPrefix("/categories").Subrouter().StrictSlash(true)
	CategoryAPI.SetCategoryRouter(a, categorySubrouter)

	// ACTIVITIES
	activitySubrouter := s.PathPrefix("/activities").Subrouter().StrictSlash(true)
	ActivityAPI.SetActivityRouter(a, activitySubrouter)
	
}