package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// AvailabilityResponse is what we expect back from Catalog.
type AvailabilityResponse struct {
	Available bool `json:"available"`
	Copies    int  `json:"copies"`
}

// CheckAvailability asks the Catalog service whether a book can be borrowed.
func CheckAvailability(bookID uint) (*AvailabilityResponse, error) {
	// Build the URL, e.g. http://localhost:8081/books/5/availability
	baseURL := os.Getenv("CATALOG_SERVICE_URL")
	url := fmt.Sprintf("%s/books/%d/availability", baseURL, bookID)

	// Make the HTTP GET request (5s timeout so we never hang forever)
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not reach catalog service: %w", err)
	}
	defer resp.Body.Close()

	// Catalog should reply 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("catalog returned status %d", resp.StatusCode)
	}

	// Read the JSON reply into our struct
	var result AvailabilityResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("could not read catalog response: %w", err)
	}

	return &result, nil
}
