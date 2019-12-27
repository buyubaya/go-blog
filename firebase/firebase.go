package firebaseAdmin


import (
	"fmt"
	"log"
	"sync"

	"golang.org/x/net/context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)


type App struct {
	Ctx context.Context
	Client *auth.Client
}


type UserInfo struct {
	UID string `json:"uid"`
	Email string `json:"email"`
	DisplayName string `json:"displayName"`
	PhotoURL string `json:"photoURL"`
}


func (fbA *App) Initialize(serviceAccountKey string) {
	fbA.Ctx = context.Background()

	opt := option.WithCredentialsFile(serviceAccountKey)
	app, err := firebase.NewApp(fbA.Ctx, nil, opt)
	if err != nil {
		log.Panic(fmt.Errorf("error initializing app: %v", err))
	}

	client, err := app.Auth(fbA.Ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}


	fbA.Client = client
}


// func (fbA *App) GetAllUsers() {

// 	// Note, behind the scenes, the Users() iterator will retrive 1000 Users at a time through the API
// 	iter := fbA.Client.Users(fbA.Ctx, "")
// 	for {
// 		user, err := iter.Next()
// 		if err == iterator.Done {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatalf("error listing users: %s\n", err)
// 		}
// 		log.Printf("read user user: %v\n", user)
// 	}

// 	// Iterating by pages 100 users at a time.
// 	// Note that using both the Next() function on an iterator and the NextPage()
// 	// on a Pager wrapping that same iterator will result in an error.
// 	pager := iterator.NewPager(fbA.Client.Users(fbA.Ctx, ""), 100, "")
// 	for {
// 		var users []*auth.ExportedUserRecord
// 		nextPageToken, err := pager.NextPage(&users)
// 		if err != nil {
// 			log.Fatalf("paging error %v\n", err)
// 		}
// 		for _, u := range users {
// 			log.Printf("read user user: %v\n", u)
// 		}
// 		if nextPageToken == "" {
// 			break
// 		}
// 	}
// }


func (fbA *App) GetUserByUID(uid string) (*UserInfo, error) {

	user, err := fbA.Client.GetUser(fbA.Ctx, uid)
	if err != nil {
		return nil, err
	}

	userInfo := &UserInfo{
		UID: user.UID,
		DisplayName: user.DisplayName,
		Email: user.Email,
		PhotoURL: user.PhotoURL,
	}

	return userInfo, nil
}


func (fbA *App) GetUserByUIDs(uids []string) (map[string]*UserInfo, error) {
	userInfosMap := make(map[string]*UserInfo)
	channel := make(chan *UserInfo)
	var wg sync.WaitGroup

	for index, _ := range uids {
		uid := uids[index]

		wg.Add(1)
		go func() {
			defer wg.Done()
			userInfo, _ := fbA.GetUserByUID(uid)

			if userInfo != nil {
				channel <- userInfo
			}
		}()
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	for userInfo := range channel {
		userInfosMap[userInfo.UID] = userInfo
	}

	return userInfosMap, nil
}


func (fbA *App) GetUserByEmail(email string) (*auth.UserRecord, error) {

	userInfo, err := fbA.Client.GetUserByEmail(fbA.Ctx, email)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}


func (fbA *App) CreateUser(user *auth.UserToCreate) (*UserInfo, error) {
	
	createdUser, err := fbA.Client.CreateUser(fbA.Ctx, user)
    if err != nil {
		return nil, err
	}
	
	userInfo := &UserInfo{
		UID: createdUser.UID,
		DisplayName: createdUser.DisplayName,
		Email: createdUser.Email,
		PhotoURL: createdUser.PhotoURL,
	}

    return userInfo, nil
}


func (fbA *App) UpdateUser(uid string, user *auth.UserToUpdate) (*UserInfo, error) {
	
	newUserInfo, err := fbA.Client.UpdateUser(fbA.Ctx, uid, user)
	if err != nil {
		return nil, err
	}

	userInfo := &UserInfo{
		UID: newUserInfo.UID,
		DisplayName: newUserInfo.DisplayName,
		Email: newUserInfo.Email,
		PhotoURL: newUserInfo.PhotoURL,
	}

    return userInfo, nil
}


func (fbA *App) VerifyIDToken(idToken string) (*auth.Token, error) {
	
	token, err := fbA.Client.VerifyIDToken(fbA.Ctx, idToken)
	if err != nil {
		return nil, err
	}

    return token, nil
}