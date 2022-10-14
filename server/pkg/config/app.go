package config

//packeage config

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "gorm.io/driver/mysql"
)

var (
	db *gorm.DB
)

func Connect() {
	d, err := gorm.Open("mysql", "vitalis:password@tcp(movies-mysql:3306)/movies?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db = d
	fmt.Println("Connected to DB")
}

func GetDb() *gorm.DB {
	return db
}
