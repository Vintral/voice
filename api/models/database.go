package models

import (
	"fmt"
	"os"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Database(retry bool) (*gorm.DB, error) {
	// Use cached value if we can
	if db != nil {
		return db, nil
	}

	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")

	dsn := DB_USER + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?charset=utf8mb4&parseTime=True&loc=Local"
	dbase, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		if retry {
			fmt.Println(err)
			panic("Failed to connect to database @ " + DB_HOST + ":" + DB_PORT)
		} else {
			time.Sleep(3 * time.Second)
			return Database(true)
		}
	}

	// Cache this to re-use next time
	db = dbase

	if err := dbase.Use(otelgorm.NewPlugin()); err != nil {
		panic(err)
	}

	sql, err := dbase.DB()
	if err != nil {
		panic(err)
	}

	sql.SetMaxOpenConns(10)
	sql.SetMaxIdleConns(3)
	sql.SetConnMaxIdleTime(5 * time.Minute)

	return dbase, err
}

func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Issue{}, &Donation{})
	fmt.Println("Ran Migrations")
}
