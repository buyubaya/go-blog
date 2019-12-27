package interfaces


type UserInfo struct {
	UID string `json:"uid"`
	Email string `json:"email"`
	DisplayName string `json:"displayName"`
	PhotoURL string `json:"photoURL"`
}