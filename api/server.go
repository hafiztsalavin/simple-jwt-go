package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"simple-jwt-go/api/routers"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()
	r.HandleFunc("/healthz", homePage)

	// Initialize routes
	routers.UserRoutes(r)

	writeTimeout, _ := strconv.Atoi(os.Getenv("WRITE_TIMEOUT_SEC"))
	readTimeout, _ := strconv.Atoi(os.Getenv("READ_TIMEOUT_SEC"))

	// Get port from env.
	server := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("localhost:%s", os.Getenv("API_PORT")),
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
	}

	fmt.Printf("Start server at localhost:%s!\n", os.Getenv("API_PORT"))
	log.Fatalln(server.ListenAndServe())
}

func homePage(w http.ResponseWriter, r *http.Request) {
	current_env := os.Getenv("ENV")
	if current_env == "" {
		current_env = "dev"
	}
	fmt.Fprintln(w, "Welcome to Golang REST - API")
	fmt.Fprintf(w, "ENV: %s", current_env)
}
