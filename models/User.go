package models


import (
	"time"
)


// MODEL
type User struct {
	FirebaseUID string `gorm:"PRIMARY_KEY" json:"firebaseUID"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`
}


func (u *User) TableName() string {
    return "h_user"
}