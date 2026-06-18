package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// 1. Fetch configurations from Environment Variables
	apiBaseURL := os.Getenv("API_BASE_URL")
	apiToken := os.Getenv("API_TOKEN")

	if apiBaseURL == "" || apiToken == "" {
		log.Fatal("API_BASE_URL and API_TOKEN environment variables must be set!")
	}

	// 2. Define the HTTP handler
	http.HandleFunc("/fetch-data", func(w http.ResponseWriter, r *http.Request) {
		// Create a new HTTP request to the external API
		req, err := http.NewRequest("GET", apiBaseURL, nil)
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		// Inject the sensitive token securely into the header
		req.Header.Add("Authorization", "Bearer "+apiToken)

		// Execute the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Failed to fetch data from external API", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Read and return the response to the user
		body, _ := io.ReadAll(resp.Body)
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})

	// 3. Start the server
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/* Understanding the Go Code:

os.Getenv(...): This is the crucial bridge between Kubernetes and your application. Instead of hardcoding the URL or the secret token in the Go code, we instruct Go to look for environment variables. Kubernetes will inject these variables into the container at runtime.

http.NewRequest(...) & req.Header.Add(...): We construct a request to the external API and append the secure token to the headers (as a Bearer token). This ensures the token is transmitted securely to the third-party service.

http.ListenAndServe(":8080", nil): This spins up a web server on port 8080 to listen for incoming traffic.

*/