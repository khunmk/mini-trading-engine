package model

import (
	"fmt"
	"log"

	"github.com/khunmk/mini-trading-engine/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	host := config.MsHost.Host
	port := config.MsHost.Port
	user := config.MsHost.User
	pass := config.MsHost.Pass
	name := config.MsHost.Name

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//
	})

	if err != nil {
		log.Fatalf("error on %v\n", err.Error())
	}

	DB = db

	db.AutoMigrate(
		&Order{},
		&Trade{},
	)
}
