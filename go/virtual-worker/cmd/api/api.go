package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Handles the Healthz endpoint
func (app *Config) HandleHealthz(w http.ResponseWriter, r *http.Request) {

	//set app port
	daprHttpPort := "5980"
	if value, ok := os.LookupEnv("DAPR_HTTP_PORT"); ok {
		daprHttpPort = value
	}
	_, err := http.Get("http://localhost:" + daprHttpPort + "/v1.0/healthz")

	if err != nil {
		app.writeError(w, err, http.StatusInternalServerError)
		os.Exit(1)
	}

	app.writeJSON(w, http.StatusOK, "Healthy")
}

// Handles the generate receipt endpoint
func (app *Config) HandleCronBinding(w http.ResponseWriter, r *http.Request) {

	const ReceiptAppId = "receipt-generation-service"

	log.Print("Calling receipt service through cron binding")
	daprHttpPort := os.Getenv("DAPR_HTTP_PORT")
	if daprHttpPort == "" {
		daprHttpPort = "5980"
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	url := fmt.Sprintf("http://localhost:%s/crontest", daprHttpPort)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("error creating request to receipt service - %s", err)
	}

	// Adding target app id as part of the header
	req.Header.Add("dapr-app-id", ReceiptAppId)

	// Invoking a service
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error calling receipt service - %s", err)
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading receipt service response body - %s", err)
	}

	app.writeJSON(w, http.StatusOK, "Cron binding invoked successfully")
}
