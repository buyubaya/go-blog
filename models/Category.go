package models


type Category struct {
	ID uint64 `gorm:"PRIMARY_KEY; AUTO_INCREMENT" json:"id"`
	Name string `gorm:"SIZE 100; NOT NULL" json:"name"`
	BackgroundColor *string `json:"backgroundColor"`
}


func (c *Category) TableName() string {
	return "h_category"
}