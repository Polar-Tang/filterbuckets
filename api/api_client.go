package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type FileInfo struct {
	URL          string `json:"url"`
	Filename     string `json:"filename"`
	FullPath     string `json:"fullPath"`
	Size         int    `json:"size"`
	LastModified int    `json:"lastModified"`
}

type ApiResponse struct {
	Files []FileInfo `json:"files"`
}

func QueryFiles(sessionCookie string, keywords []string, extensions []string) ([]FileInfo, error) {
	apiURL := "https://buckets.grayhatwarfare.com/api/v2/files"

	// Build query parameters
	params := url.Values{}
	params.Set("keywords", joinKeywords(keywords))
	params.Set("extensions", joinKeywords(extensions))
	params.Set("limit", "1000")

	// Build the full URL
	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	// Create HTTP client with timeout
	client := &http.Client{Timeout: 30 * time.Second}

	// Create the request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers (including the authorization token)
	req.Header.Set("Cookie", "Cookie: "+sessionCookie)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for successful status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	// Parse the JSON response
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return apiResponse.Files, nil
}

// Helper function to join keywords into a single string
func joinKeywords(keywords []string) string {
	return url.QueryEscape(strings.Join(keywords, " "))
}
