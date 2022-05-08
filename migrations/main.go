package migrations

import (
	"fmt"
	"log"
	"simple-jwt-go/api/database"
	"simple-jwt-go/api/models"
)

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

// MigrateModels used to migrate table of models that need to be created in the database by dropping existing tables then re-create it
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
