package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"log"
	"fmt"
	"sync"

	"challenge3/models"
)

var (
	lock = &sync.Mutex{}
	connectionSingleton *gorm.DB
)

func GetDatabase() *gorm.DB {
	if connectionSingleton == nil {
		lock.Lock()
		defer lock.Unlock()

		if connectionSingleton == nil {
			dbname := "Challenge3"
			db := "postgres"
			dbpassword := "Cuongnguyen2001"
			dburl := "postgres://postgres:" + dbpassword + "@challenge3/" + dbname + "?sslmode=disable"

			connection, err := gorm.Open(db, dburl)

			if err != nil {
				fmt.Println(err.Error())
				log.Fatalln("Wrong database url")
			}

			sqldb := connection.DB()
			err = sqldb.Ping()

			if err != nil {
				log.Fatalln("Database is connectedt")
			}
			// log.Fatalln("Database is connected")
			fmt.Println("Database is connected")
			return connection
		}

		fmt.Println("Database is already connected")
	}

	fmt.Println("Database is already connected")

	return connectionSingleton
}

func InitMigration() {
	connection := GetDatabase()
	// defer CloseDatabase()

	connection.AutoMigrate(models.User{})
	connection.AutoMigrate(models.Post{})
	connection.AutoMigrate(models.Role{})
	connection.AutoMigrate(models.CheckPass{})
	connection.AutoMigrate(models.PasswordResetToken{})

}

func CloseDatabase() {
	sqldb := connectionSingleton.DB()
	sqldb.Close()
}