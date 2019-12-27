package models


import (
	"time"
)


type Activity struct {
	ID uint64 `json:"id"`
	Action string `json:"action"`
	Changelog string `json:"changelog"`
	CreatedAt time.Time `json:"creaatedAt"`
	CreatedBy string `json:"createdBy"`
}


func (a *Activity) TableName() string {
	return "h_activity"
}