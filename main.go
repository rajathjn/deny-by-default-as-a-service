package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rajathjn/deny-by-default-as-a-service/cmd"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := "0.0.0.0:" + port

	// Work around for Health check mode for Docker healthcheck (scratch image has no curl/wget)
	if len(os.Args) > 1 && os.Args[1] == "-health" {
		response, err := http.Get("http://localhost:" + port + "/health")
		if err != nil || response.StatusCode != http.StatusOK {
			fmt.Println("unhealthy")
			os.Exit(1)
		}
		fmt.Println("healthy")
		os.Exit(0)
	}

	cmd.Server(address)
}
