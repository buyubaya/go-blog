package models


type ObjectFieldActivity struct {
	ObjectID string `json:"objectID"`
	ActivityID uint64 `json:"activityID"`
}


func (a *ObjectFieldActivity) TableName() string {
	return "h_object_field_activity"
}