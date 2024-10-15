package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const ReceiptAppId = "receipt-generation-service"

type Config struct{}

func main() {
	daprHttpPort := os.Getenv("DAPR_HTTP_PORT")
	if daprHttpPort == "" {
		daprHttpPort = "5980"
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	url := fmt.Sprintf("http://localhost:%s/ready", daprHttpPort)
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

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading receipt service response body - %s", err)
	}

	fmt.Println(string(b))

}
