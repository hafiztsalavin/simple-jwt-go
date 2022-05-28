package seeders

import (
	"simple-jwt-go/api/database"
	"simple-jwt-go/api/models"
	"simple-jwt-go/api/utils"
)

func UserSeeder() error {

	// connect db
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	pass, _ := utils.HashPassword("123pass")
	user := models.User{
		Username: "Alvin",
		Email:    "alvintest@gmail.com",
		Password: pass,
	}

	if result := db.Create(&user); result.Error == nil {
		return result.Error
	}

	return nil
}
