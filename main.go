package main

import (
	"log"
	"os"

	"simple-jwt-go/api"
	"simple-jwt-go/migrations"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error cant read env, %v", err)
	} else {
		if len(os.Args) > 2 && os.Args[2] == "migrate" {
			if err := migrations.MigrateModels(); err != nil {
				log.Fatalf("Error when migrate models, %v", err)
			}
		} else {
			api.Run()
		}
	}
}
