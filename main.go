package main

// IMPORTS
import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"pdf_greyhat_go/api"
	"pdf_greyhat_go/download"
)

func main() {
	// Initialize session and keywords
	sessionCookie := "01931a3ff4929fa0e8d8c93ba9dac24c"
	keywords := []string{"deploy"}
	extensions := map[string][]string{
		"json": {"token", "credentials", "password", "key", "secret", "id", "name"},
	}
	for _, keyword := range keywords {
		outputFile := fmt.Sprintf("results-%s.json", keyword)
		fmt.Printf("Searching for files with keyword: %s\n", keyword)
		// Query files THROUGH THE API
		files, err := api.QueryFiles(sessionCookie, []string{keyword}, extensions)
		// generic error handling
		if err != nil {
			log.Printf("Failed to query files: %v", err)
			continue
		}

		// Create a semaphore for concurrent downloads
		var wg sync.WaitGroup

		// Initialize results
		results := make([]map[string]interface{}, 0)
		// RESULTS:
		/*
				{"Filename": "file1.pdf", "URL": "http://example.com/file1", "Matches": 10},
			    {"Filename": "file2.pdf", "URL": "http://example.com/file2", "Matches": 5},
		*/

		mutex := &sync.Mutex{}

		// Set the concurrency limit
		concurrencyLimit := 6
		// use semaphore var to set a maximum number of concurrent goroutines
		semaphore := make(chan struct{}, concurrencyLimit)

		// Creates a timer that triggers every 60 seconds.
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		// iterate just on the values of the files
		// go routine
		// HANDLE THE GO ROUTINE CONCURRENCY

		// ensure to goes periodicly saving it's making to avoid big lost if the process is interrupted
		go func() {
			// A channel that emits a signal every time the ticker fires.
			for range ticker.C {
				// save the file periodically
				mutex.Lock()
				saveResults(results, outputFile)
				fmt.Printf("Result added: %+v\n", results) // Add this line for debugging
				// MUTEX write the file but priventing race conditions
				mutex.Unlock()
			}
		}()

		for _, fileInfo := range files {
			// increment the wait counter
			wg.Add(1)
			go func(file api.FileInfo) {
				// DECREMENT the wait routine when it's done
				defer wg.Done()
				// send an empty struct into the sempahore channel
				semaphore <- struct{}{} // Acquire a semaphore slot
				// semaphoro green!
				defer func() { <-semaphore }() // Release slot after processing

				// fmt.Printf("Found file: %s (URL: %s, Size: %d bytes)\n", file.Filename, file.URL, file.Size)

				result := download.ProcessFile(file, extensions) // redefine result
				// redefine the results with the function proces file
				if result != nil {
					// append the result (no overwrite)
					mutex.Lock()
					results = append(results, result)

					// MUTEX write the file but priventing race conditions
					mutex.Unlock()
				}
			}(fileInfo)
		}
		// The file info is a struct
		/* type FileInfo struct {
			URL      string
			Filename string
			Size     int
		} */
		// The results are saved as JSON in results.json, after the whole fucking process ends:
		saveResults(results, outputFile)
		wg.Wait() // Wait for all goroutines to complete
	}
}

func saveResults(results []map[string]interface{}, outputFile string) {
	fmt.Printf("Saving %d results...\n", len(results)) // Debug log

	file, err := os.Create(outputFile) // Create (or overwrite) results.json
	if err != nil {
		log.Printf("Failed to create output file: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Add indentation for readability

	if err := encoder.Encode(results); err != nil {
		log.Printf("Failed to write results: %v", err)
	} else {
		fmt.Printf("Results saved to %s\n", outputFile)
	}
}
