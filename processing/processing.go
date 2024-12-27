package processing

// IMPORTS
import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Polar-Tang/filterbuckets/api"
	"github.com/Polar-Tang/filterbuckets/download"

	"github.com/fatih/color"
)

var (
	wg    sync.WaitGroup
	mutex sync.Mutex
	err   error
	files []api.FileInfo
)

// keywords []string, extensions map[string][]string

func ProcessFiles(keywords []string, extensions map[string][]string, bucketFile string) {

	sessionCookie := "54e7fe8c2aa1dd504b9be39fa3466f10"
	results := make([]map[string]interface{}, 0)

	ticker := time.NewTicker(300 * time.Second)
	defer ticker.Stop()
	tickerColor := color.New(color.FgBlue).PrintlnFunc()

	go func() {
		for range ticker.C {
			tickerColor("Periodic save: Saving current results...")
			mutex.Lock()
			err := SaveResults(results, "generic_results.json")
			mutex.Unlock()
			if err != nil {
				log.Printf("Error during periodic save: %v", err)
			} else {
				tickerColor("Periodic save complete.")
			}
		}
	}()

	// ---------------------------------------------------------------

	if len(keywords) == 0 {
		fmt.Println("Processing files without keywords...")
		files, err = api.QueryFiles(sessionCookie, []string{}, extensions, bucketFile)
		fmt.Println(bucketFile)
		ProcessFileForKeyword("", extensions, sessionCookie, results)
	} else {
		for _, keyword := range keywords {

			cleankeyword := strings.TrimSpace(keyword)
			cleankeyword = strings.Trim(cleankeyword, `"`)
			cleankeyword = strings.TrimRight(cleankeyword, `",`)

			outputFile := createOutputFile(cleankeyword)
			if err != nil {
				fmt.Printf("Failed to create output file: %v\n", err)
				continue
			}
			fmt.Printf("Searching for files with keyword: %s\n", cleankeyword)

			// -------------------------------------------------------------

			maxRetries := 3
			for retries := 0; retries < maxRetries; retries++ {
				files, err = api.QueryFiles(sessionCookie, []string{cleankeyword}, extensions, bucketFile)
				if err == nil {
					break
				}
				log.Printf("Retry %d/%d for keyword '%s' failed: %v", retries+1, maxRetries, cleankeyword, err)
				time.Sleep(2 * time.Second)
			}
			if err != nil {
				log.Printf("All retries failed for keyword '%s'\n", cleankeyword)
				return
			}
			ProcessFileForKeyword(keyword, extensions, sessionCookie, results)
			mutex.Lock()
			err = SaveResults(results, outputFile)
			if err != nil {
				log.Printf("Error saving final results for keyword '%s': %v", keyword, err)
			}
			mutex.Unlock()
		}
		// -------------------------------------------------------------

		wg.Wait()

	}

}

func SaveResults(results []map[string]interface{}, outputFile string) error {
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

func ProcessFileForKeyword(keyword string, extensions map[string][]string, sessionCookie string, results []map[string]interface{}) {
	errorschan := make(chan error, len(files))
	concurrencyLimit := 6
	semaphore := make(chan struct{}, concurrencyLimit)
	outputFile := createOutputFile(keyword)

	done := make(chan struct{})
	defer close(done)

	for _, fileInfo := range files {

		wg.Add(1)

		if fileInfo.Size > 50*1024*1024 {
			errorschan <- fmt.Errorf("skipping large file: %s", fileInfo.Filename)
			continue
		}

		// 💚💚💚💚💚💚💚💚💚💚💚
		processingColor := color.New(color.FgGreen).PrintlnFunc()
		go func(file api.FileInfo) {
			semaphore <- struct{}{}
			processingColor("Starting a goroutine...")

			defer func() {
				processingColor("Exiting goroutine...")
			}()
			defer wg.Done()

			start := time.Now()
			result := download.ProcessFile(file, extensions)
			processingColor("File processed in: ", time.Since(start))
			<-semaphore
			if result != nil {
				processingColor("Locking mutex...\n")
				mutex.Lock()
				results = append(results, result)
				processingColor("Unocking mutex...\n")
				mutex.Unlock()
			}
		}(fileInfo)
	}
	wg.Wait()
	// TICKER

	// 💛💛💛💛💛💛💛💛💛💛💛
	selectingColor := color.New(color.FgYellow).PrintlnFunc()
	go func() {
		selectingColor("Starting a goroutine...")
		defer selectingColor("Exiting goroutine...")
		for {
			select {
			case err, ok := <-errorschan:
				if !ok {
					break
				}
				fmt.Printf("Error: %v\n", err)
			case _, ok := <-done:
				if !ok {
					break
				}
				selectingColor("All files processed")
				return
			}
		}
	}()
	fmt.Printf("Saving results for keyword: %s\n", keyword)
	mutex.Lock()
	err := SaveResults(results, outputFile)
	mutex.Unlock()
	if err != nil {
		log.Printf("Error saving results for keyword '%s': %v", keyword, err)
	}
}

// ---------------------------------------------------------------
func createOutputFile(keyword string) string {
	if keyword == "" {
		genericFilename := "results.json"
		return genericFilename
	}
	filename := fmt.Sprintf("results-%s.json", keyword)
	dir, err := os.Open(".")
	if err != nil {
		return ""
	}
	defer dir.Close()
	fmt.Printf("Checking for existing file: %s", filename)
	var acc int
	names, err := dir.Readdirnames(-1)
	if err != nil && err != io.EOF {
		return ""
	}

	for _, name := range names {
		if name == filename || name == fmt.Sprintf("results-%s-%d.json", keyword, acc) {
			acc++
		}
	}

	if acc > 0 {
		filename = fmt.Sprintf("results-%s-%d.json", keyword, acc)
	}

	return filename
}
