package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Struct to hold file information from the API response
type FileInfo struct {
	URL        string `json:"url"`
	Filename   string `json:"filename"`
	FullPath   string `json:"fullPath"`
	Size       int    `json:"size"`
	LastModified int  `json:"lastModified"`
}

// Struct to parse the entire API response
type ApiResponse struct {
	Files []FileInfo `json:"files"`
}

// QueryFiles fetches files from the Greyhat API
func QueryFiles(sessionCookie string, keywords []string, extensions []string) ([]FileInfo, error) {
	apiURL := "https://buckets.grayhatwarfare.com/api/v2/files"

	// Build query parameters
	params := url.Values{}
	params.Set("keywords", url.QueryEscape(joinKeywords(keywords)))
	params.Set("extensions", url.QueryEscape(joinKeywords(extensions)))
	params.Set("limit", "1000") // Set limit to 100 (you can customize this)

	// Build the full URL
	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	// Create HTTP client with timeout
	client := &http.Client{Timeout: 30 * time.Second}

	// Create the request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	// Add headers (including the authorization token)
	req.Header.Set("Cookie", "Cookie: "+sessionCookie)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check for successful status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	// Parse the JSON response
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse.Files, nil
}

// Helper function to join keywords into a single string
func joinKeywords(keywords []string) string {
	return url.QueryEscape(joinSlice(keywords))
}

// Helper function to join a slice of strings with spaces
func joinSlice(slice []string) string {
	return fmt.Sprintf("%s", url.QueryEscape(slice[0]))
}



