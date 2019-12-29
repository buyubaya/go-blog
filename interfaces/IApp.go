package interfaces


import (
	"github.com/jinzhu/gorm"

	firebaseAdmin "github.com/buyubaya/go-blog/firebase"
	"github.com/go-redis/redis/v7"
)


type (
	// FirebaseAdminApp interface {
	// 	VerifyIDToken(idToken string) (*auth.Token, error)
	// 	GetUserByUID(uid string) (*UserInfo, error)
	// 	GetUserByUIDs(uids []string) (map[string]*UserInfo, error)
	// 	GetUserByEmail(email string) (*auth.UserRecord, error)
	// 	CreateUser(user *auth.UserToCreate) (*UserInfo, error)
	// }


	IApp interface {
		GetDB() *gorm.DB
		GetFirebaseApp() *firebaseAdmin.App
		GetRedisClient() *redis.Client
	}
)