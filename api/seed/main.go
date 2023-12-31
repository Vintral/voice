package main

import (
	"api/models"
	"fmt"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db, err := models.Database(false)
	if err != nil {
		panic(err)
	}

	db.Exec("DROP TABLE donations")
	db.Exec("DROP TABLE users")
	db.Exec("DROP TABLE issues")

	models.RunMigrations(db)

	fmt.Println("Seeding users")
	db.Create(&models.User{Email: "jane.doe@email.com", Admin: true})
	db.Create(&models.User{Email: "john.doe@email.com"})
	db.Create(&models.User{Email: "another@email.com"})
	db.Create(&models.User{Email: "joe@email.com"})

	fmt.Println("Seeding issues")
	for i := 1; i < 5; i++ {
		db.Create(&models.Issue{
			GUID:    uuid.New(),
			Title:   fmt.Sprintf("Issue %d", i),
			Summary: fmt.Sprintf("Summary about Issue %d would go here.  Lorem ipsum", i),
			Pool:    0.00,
		})
	}
}
