package main

import (
	"log"
	"net/http"
	"os"
)

type Config struct {
}

func main() {

	//set app port
	appPort := "5900"
	if value, ok := os.LookupEnv("APP_PORT"); ok {
		appPort = value
	}

	app := Config{}

	log.Printf("Starting the application on port %s\n", appPort)

	// create a new server
	srv := &http.Server{
		Addr:    ":" + appPort,
		Handler: app.routes(),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
