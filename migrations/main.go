package migrations

import (
	"fmt"
	"log"
	"simple-jwt-go/api/database"
	"simple-jwt-go/api/models"
)

// drop tables
func DropTables() error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	if err := db.Migrator().DropTable("users"); err != nil {
		return err
	}

	fmt.Println("Table(s) dropped")

	return nil
}

// make a new table based on models
func MigrateModels() error {
	if err := DropTables(); err != nil {
		return err
	}

	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("Model(s) migrated")

	return nil
}
