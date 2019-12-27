package controllers


import (
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/gorilla/mux"
	firebaseAdmin "github.com/buyubaya/go-blog/firebase"

	// "github.com/buyubaya/go-blog/api/interfaces"

	"github.com/buyubaya/go-blog/config"
	// "github.com/buyubaya/go-blog/controllers"
	"github.com/buyubaya/go-blog/models"
	// "github.com/buyubaya/go-blog/api/middlewares"
	postHandlers "github.com/buyubaya/go-blog/handlers/post"
)


type App struct {
	DB *gorm.DB
	Router *mux.Router
	FirebaseApp *firebaseAdmin.App
}


// Initialize initializes the app with predefined configuration
func (a *App) Initialize(config *config.Config) {
	
	// DB
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
		config.DB.Charset)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	if err != nil {
		log.Fatal("Could not connect database", err.Error())
	}

	a.DB = models.DBMigrate(db)


	// ROUTER
	a.Router = mux.NewRouter()
	a.setRouters()


	// FIREBASE
	firebaseApp := &firebaseAdmin.App{}
	firebaseApp.Initialize(config.Firebase.ServiceAccountKey)
	a.FirebaseApp = firebaseApp
}


// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}


// METHODS
func (a *App) GET(path string, f http.HandlerFunc) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

func (a *App) POST(path string, f http.HandlerFunc) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

func (a *App) PUT(path string, f http.HandlerFunc) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

func (a *App) DELETE(path string, f http.HandlerFunc) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}


/******************** API ROUTER ********************/
func (a *App) setRouters() {
	// BLOGS
	a.GET("/posts", postHandlers.GetAllPosts)
}