package util

import (
	"fmt"
	"log"

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
	connectionString := "host=localhost user=postgres dbname=UserServiceDB sslmode=disable password=root port=5432"
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
