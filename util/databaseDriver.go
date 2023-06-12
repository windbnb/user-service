package util

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	model "github.com/windbnb/user-service/model"
)

var (
	users = []model.User{
		{Email: "host@email.com", Username: "ivica98", Password: "host", Name: "Ivica", Surname: "Roganovic", Address: "Maksima Gorkog 17a, Novi Sad", Role: model.HOST},
		{Email: "guest@email.com", Username: "makulica", Password: "guest", Name: "Jovana", Surname: "Mustur", Address: "Dr Svetislava Kasapinovica 22, Novi Sad",Role: model.GUEST},
	}
)

func ConnectToDatabase() *gorm.DB {
	host, hostFound := os.LookupEnv("DATABASE_HOST")
	if !hostFound {
		host = "localhost"
	}
	user, userFound := os.LookupEnv("DATABASE_USER")
	if !userFound {
		user = "postgres"
	}
	password, passwordFound := os.LookupEnv("DATABASE_PASSWORD")
	if !passwordFound {
		password = "root"
	}

	connectionString := "host=" + host + " user=" + user + " dbname=UserServiceDB sslmode=disable password=" + password + " port=5432"
	dialect := "postgres"

	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection to DB successfull.")
	}

	db.DropTable("users")
	db.DropTable("user_deletion_events")
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.UserDeletionEvent{})

	for _, user := range users {
		db.Create(&user)
	}

	return db
}
