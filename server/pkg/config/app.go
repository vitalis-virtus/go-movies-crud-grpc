package config

//package config

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "gorm.io/driver/mysql"
)

var (
	db *gorm.DB
)

func Connect() {
	d, err := gorm.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/movies?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db = d
	fmt.Println("Connected to DB")
}

func GetDb() *gorm.DB {
	return db
}
