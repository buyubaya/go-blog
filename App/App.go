package App


import (
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	firebaseAdmin "github.com/buyubaya/go-blog/firebase"
	"github.com/go-redis/redis/v7"

	"github.com/buyubaya/go-blog/config"
	"github.com/buyubaya/go-blog/models"
	API "github.com/buyubaya/go-blog/api"
	UserAPI "github.com/buyubaya/go-blog/api/user"
	"github.com/buyubaya/go-blog/handlers"
	redisClient "github.com/buyubaya/go-blog/redisClient"
)


type App struct {
	DB *gorm.DB
	Router *mux.Router
	FirebaseApp *firebaseAdmin.App
	RedisClient *redis.Client
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


	// REDIS CLIENT
	redisClient := redisClient.Initialize(config.Redis)
	a.RedisClient = redisClient
}


/******************** GETTER ********************/
func (a *App) GetDB() *gorm.DB {
	return a.DB
}

func (a *App) GetFirebaseApp() *firebaseAdmin.App {
	return a.FirebaseApp
}

func (a *App) GetRedisClient() *redis.Client {
	return a.RedisClient
}


// Run the app on it's router
func (a *App) Run(host string) {
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders: []string{"Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
		Debug: true,
	}).Handler(a.Router)
	
	log.Fatal(http.ListenAndServe(host, handler))
}


/******************** API ROUTER ********************/
func (a *App) setRouters() {

	// LOGIN
	a.Router.HandleFunc("/login", handlers.Login(a)).Methods("POST")

	// REGISTER
	a.Router.HandleFunc("/register", UserAPI.CreateUser(a)).Methods("POST")

	// APIs
	s := a.Router.PathPrefix("/api").Subrouter().StrictSlash(true)
	API.SetAPIRouter(a, s)
}