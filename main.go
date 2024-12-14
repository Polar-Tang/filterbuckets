package main

// IMPORTS
import (
	"encoding/json"
	"fmt"
	"io"
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
	keywords := []string{
		"manifest",
		"mbank",
	}
	extensions := map[string][]string{
		"json": {".*mbank.*"},
		"xml":  {".*mbank.*"},
		"pdf":  {".*mbank.*"},
		"php":  {".*mbank.*"},
		"js":   {".*mbank.*"},
		"py":   {".*mbank.*"},
		"java": {".*mbank.*"},
		"go":   {".*mbank.*"},
		"txt":  {".*mbank.*"},
		"html": {".*mbank.*"},
	}

	createOutputFile := func(keyword string) (string, error) {
		filename := fmt.Sprintf("results-%son", keyword)
		dir, err := os.Open(".")
		if err != nil {
			return "", fmt.Errorf("failed opening the directory: %w", err)
		}
		defer dir.Close()

		var acc int
		names, err := dir.Readdirnames(-1)
		if err != nil && err != io.EOF { // EOF means end of directory
			return "", fmt.Errorf("error reading directory: %w", err)
		}

		for _, name := range names {
			if name == filename || name == fmt.Sprintf("results-%s-%don", keyword, acc) {
				acc++
			}
		}

		if acc > 0 {
			filename = fmt.Sprintf("results-%s-%don", keyword, acc)
		}

		return filename, nil
	}

	for _, keyword := range keywords {
		outputFile, err := createOutputFile(keyword)
		if err != nil {
			fmt.Printf("Failed to create output file: %v\n", err)
			continue
		}
		fmt.Printf("Searching for files with keyword: %s\n", keyword)
		var files []api.FileInfo
		maxRetries := 3
		for retries := 0; retries < maxRetries; retries++ {
			files, err = api.QueryFiles(sessionCookie, []string{keyword}, extensions)
			if err == nil {
				break // Exit the retry loop if successful
			}
			log.Printf("Retry %d/%d for keyword '%s' failed: %v", retries+1, maxRetries, keyword, err)
			time.Sleep(2 * time.Second)
		}
		if err != nil {
			log.Printf("All retries failed for keyword '%s'\n", keyword)
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

		// ensure to goes periodicly saving it's making to avoid big lost if the process is interrupted
		go func() {
			// A channel that emits a signal every time the ticker fires.
			for range ticker.C {
				// save the file periodically
				mutex.Lock()
				err := saveResults(results, outputFile)
				if err != nil {
					log.Printf("Error saving periodic results for keyword '%s': %v", keyword, err)
				}
				// fmt.Printf("Result added: %+v\n", results) // Add this line for debugging
				// MUTEX write the file but priventing race conditions
				mutex.Unlock()
			}
		}()

		for _, fileInfo := range files {

			fmt.Println("Processing file:", fileInfo.Filename)
			if fileInfo.Size > 50*1024*1024 { // Skip files larger than 50 MB
				fmt.Printf("Skipping large file: %s\n", fileInfo.Filename)
				continue
			}

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
		// The results are saved as JSON in resultson, after the whole fucking process ends:
		wg.Wait() // Wait for all goroutines to complete
		mutex.Lock()
		err = saveResults(results, outputFile)
		if err != nil {
			log.Printf("Error saving final results for keyword '%s': %v", keyword, err)
		}
		mutex.Unlock()
	}
}

func saveResults(results []map[string]interface{}, outputFile string) error {
	fmt.Printf("Saving %d results...\n", len(results)) // Debug log

	file, err := os.Create(outputFile) // Create (or overwrite) resultson
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %w", outputFile, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Add indentation for readability

	if err := encoder.Encode(results); err != nil {
		return fmt.Errorf("failed to write JSON to file '%s': %w", outputFile, err)
	}
	fmt.Printf("Results saved to %s\n", outputFile)
	return nil
}
