package main

import (
	"log"
	"os"

	"simple-jwt-go/api"
	"simple-jwt-go/migrations"
	"simple-jwt-go/seeders"

	"github.com/joho/godotenv"
)

func main() {
	// Find and read the config file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error cant read env, %v", err)
	} else {
		if len(os.Args) > 2 && os.Args[2] == "migrate" {
			if err := migrations.MigrateModels(); err != nil {
				log.Fatalf("Error when migrate models, %v", err)
			}
		} else if len(os.Args) > 2 && os.Args[2] == "seed" {
			if err := seeders.UserSeeder(); err != nil {
				log.Fatalf("Error when migrate models, %v", err)
			}
		} else {
			api.Run()
		}
	}
}
