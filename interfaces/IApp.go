package interfaces


import (
	"github.com/jinzhu/gorm"

	firebaseAdmin "github.com/buyubaya/go-blog/firebase"
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
	}
)