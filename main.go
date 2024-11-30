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

const outputFile = "results.json"

func main() {
	// Initialize session and keywords
	sessionCookie := "__stripe_mid=af2965ba-4f1d-4bf0-a073-1aa7b3987d61ccd651; _gid=GA1.2.1907481890.1732978699; _gat_gtag_UA_121795267_1=1; SFSESSID=eu2ptgedduojcnvtch4jekfsb1; __stripe_sid=db426450-2cff-440b-9532-e088d3707d69933751; _ga=GA1.2.1930062620.1731774259; _ga_QGK3VF4QHK=GS1.1.1732995391.49.1.1732995413.0.0.0"
	keywords := []string{"mercado libre"}
	extensions := []string{"pdf"}
	pdfKeywords := []string{"mercado libre"}

	// Query files THROUGH THE API
	files, err := api.QueryFiles(sessionCookie, keywords, extensions)
	// generic error handling
	if err != nil {
		log.Fatalf("Failed to query files: %v", err)
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
	//
	semaphore := make(chan struct{}, concurrencyLimit)

	ticker := time.NewTicker(30 * time.Second) // Save every 10 seconds
	defer ticker.Stop()
	// iterate just on the values of the files
	// go routine
	// HANDLE THE GO ROUTINE CONCURRENCY

	// ensure to goes periodicly saving it's making to avoid big lost if the process is interrupted
	go func() {
		// wait for the ticker to expire
		for range ticker.C {
			// save the file periodically
			mutex.Lock()
			saveResults(results)
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

			result := download.ProcessFile(file, pdfKeywords) // redefine result
			// redefine the results with the function proces file
			if result != nil {
				// append the result (no overwrite)
				mutex.Lock()
				results = append(results, result)
				// The results are saved as JSON in results.json, after the whole fucking process ends:
				saveResults(results)
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

	wg.Wait() // Wait for all goroutines to complete

	// Final save to ensure all results are written
	saveResults(results)

}

func saveResults(results []map[string]interface{}) {
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
