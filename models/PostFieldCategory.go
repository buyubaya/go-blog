package models


type PostFieldCategory struct {
	PostID uint64 `gorm:"PRIMARY_KEY" json:"postID"`
	CategoryID uint64 `gorm:"PRIMARY_KEY" json:"categoryID"`
}


func (pfc *PostFieldCategory) TableName() string {
	return "h_post_field_category"
}