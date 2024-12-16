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
	sessionCookie := "54e7fe8c2aa1dd504b9be39fa3466f10"
	keywords := []string{
		"mbank",
	}
	extensions := map[string][]string{
		"go": {".*mbank.*"},
	}

	createOutputFile := func(keyword string) (string, error) {
		filename := fmt.Sprintf("results-%s.json", keyword)
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
			if name == filename || name == fmt.Sprintf("results-%s-%d.json", keyword, acc) {
				acc++
			}
		}

		if acc > 0 {
			filename = fmt.Sprintf("results-%s-%d.json", keyword, acc)
		}

		return filename, nil
	}

	// ---------------------------------------------------------------
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
				break
			}
			log.Printf("Retry %d/%d for keyword '%s' failed: %v", retries+1, maxRetries, keyword, err)
			time.Sleep(2 * time.Second)
		}
		if err != nil {
			log.Printf("All retries failed for keyword '%s'\n", keyword)
			continue
		}

		var wg sync.WaitGroup
		results := make([]map[string]interface{}, 0)
		var mutex sync.Mutex

		concurrencyLimit := 6
		semaphore := make(chan struct{}, concurrencyLimit)
		errorschan := make(chan error, len(files))

		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		go func() {
			for range ticker.C {
				mutex.Lock()
				err := saveResults(results, outputFile)
				if err != nil {
					errorschan <- fmt.Errorf("error saving periodic results for keyword '%s': %v", keyword, err)
				}
				mutex.Unlock()
			}
		}()
		done := make(chan struct{})
		defer close(done)
		for _, fileInfo := range files {
			wg.Add(1)

			fmt.Println("Processing file:", fileInfo.Filename)
			if fileInfo.Size > 50*1024*1024 {
				errorschan <- fmt.Errorf("skipping large file: %s", fileInfo.Filename)
				continue
			}

			go func(file api.FileInfo) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				result := download.ProcessFile(file, extensions)
				if result != nil {
					mutex.Lock()
					results = append(results, result)
					mutex.Unlock()
					if err != nil {
						errorschan <- fmt.Errorf("failed to open file '%s' for writing: %v", outputFile, err)
						return
					}
				} else {
					errorschan <- fmt.Errorf("processing failed for file: %s", file.URL)
				}
			}(fileInfo)
		}

		go func() {
			for err := range errorschan {
				fmt.Printf("Error: %v\n", err)
			}
			close(errorschan) // Close the error channel after all errors are collected
		}()

		go func() {
			for {
				select {
				case err := <-errorschan:
					fmt.Printf("Error: %v\n", err)
				case <-done:
					fmt.Println("All files processed")
					return
				}
			}
		}()
		wg.Wait()
		mutex.Lock()
		err = saveResults(results, outputFile)
		if err != nil {
			log.Printf("Error saving final results for keyword '%s': %v", keyword, err)
		}
		mutex.Unlock()
	}
}

func saveResults(results []map[string]interface{}, outputFile string) error {
	fmt.Println("Saving results...")

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %w", outputFile, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		return fmt.Errorf("failed to write JSON to file '%s': %w", outputFile, err)
	}

	fmt.Printf("Results saved to %s\n", outputFile)
	return nil
}
