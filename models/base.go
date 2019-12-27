package models


import (
	"github.com/jinzhu/gorm"
)


// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Post{}, &User{})
	return db
}