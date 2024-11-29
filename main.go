package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"pdf_processing_go/api"
	"pdf_processing_go/download"
)

const outputFile = "results.json"

func main() {
	// Initialize session and keywords
	sessionCookie := "__stripe_mid=af2965ba-4f1d-4bf0-a073-1aa7b3987d61ccd651; _gid=GA1.2.1776727429.1732643750; _gat_gtag_UA_121795267_1=1; SFSESSID=a6acupt6fcv5griv6r5tl7v6se; _ga_QGK3VF4QHK=GS1.1.1732729715.42.1.1732729727.0.0.0; _ga=GA1.1.1930062620.1731774259"
	keywords := []string{"mercado libre"}
	extensions := []string{"pdf"}
	pdfKeywords := []string{"mercado libre"}

	// Query files from API
	files, err := api.QueryFiles(sessionCookie, keywords, extensions)
	if err != nil {
		log.Fatalf("Failed to query files: %v", err)
	}

	// Create a semaphore for concurrent downloads
	var wg sync.WaitGroup
	results := make([]map[string]interface{}, 0)
	mutex := &sync.Mutex{}

	concurrencyLimit := 5
	semaphore := make(chan struct{}, concurrencyLimit)

	for _, fileInfo := range files {
		wg.Add(1)
		go func(file api.FileInfo) {
			defer wg.Done()
			semaphore <- struct{}{} // Acquire a semaphore slot
			defer func() { <-semaphore }() // Release slot after processing

			fmt.Printf("Found file: %s (URL: %s, Size: %d bytes)\n", file.Filename, file.URL, file.Size)
			result := download.ProcessFile(file, pdfKeywords)
			if result != nil {
				mutex.Lock()
				results = append(results, result)
				mutex.Unlock()
			}
		}(fileInfo)
	}

	wg.Wait() // Wait for all goroutines to complete
	saveResults(results)
}

func saveResults(results []map[string]interface{}) {
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()

	json.NewEncoder(file).Encode(results)
	fmt.Printf("Results saved to %s\n", outputFile)
}

