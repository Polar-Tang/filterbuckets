package processing

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Polar-Tang/filterbuckets/api"
	"github.com/Polar-Tang/filterbuckets/download"

	"github.com/fatih/color"
)

var (
	wg           sync.WaitGroup
	mutex        sync.Mutex
	err          error
	files        []api.FileInfo
	fileJSONName string
)

func ProcessFiles(keywords []string, extensions map[string][]string, bucketFile string, concurrencyLimit int) {

	sessionCookieFile := "./sessionCookie"

	sessionCookie, err := readSessionCookie(sessionCookieFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	results := make([]map[string]interface{}, 0)

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	tickerColor := color.New(color.FgBlue).PrintlnFunc()

	go func() {
		for range ticker.C {
			tickerColor("Periodic save: Saving current results...")
			mutex.Lock()
			err := SaveResults(results, fileJSONName)
			mutex.Unlock()
			if err != nil {
				log.Printf("Error during periodic save: %v", err)
			}
		}
	}()
	if len(keywords) == 0 {
		files, err = api.QueryFiles(sessionCookie, []string{}, extensions, bucketFile)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		ProcessFileForKeyword("", extensions, sessionCookie, results, concurrencyLimit)

	} else {
		for _, keyword := range keywords {
			cleanKeyword := strings.TrimSpace(keyword)
			cleanKeyword = strings.Trim(cleanKeyword, `"`)
			cleanKeyword = strings.TrimRight(cleanKeyword, `",`)

			fmt.Println("Processing files with keyword:", cleanKeyword)
			fileJSONName = fmt.Sprintf("results-%s.json", cleanKeyword)

			var acc int
			for {
				if _, err := os.Stat(fileJSONName); err == nil {
					acc++
					fileJSONName = fmt.Sprintf("results-%s-%d.json", cleanKeyword, acc)
					fmt.Printf("File already exists, creating new name: %s\n", fileJSONName)
					break
				} else if os.IsNotExist(err) {
					fmt.Printf("Creating: %s\n", fileJSONName)
					break
				} else {
					fmt.Printf("Error checking file: %v\n", err)
					return
				}
			}

			maxRetries := 3
			for retries := 0; retries < maxRetries; retries++ {
				files, err = api.QueryFiles(sessionCookie, []string{cleanKeyword}, extensions, bucketFile)
				if err == nil {
					break
				}
				log.Printf("Retry %d/%d for keyword '%s' failed: %v", retries+1, maxRetries, cleanKeyword, err)
				time.Sleep(2 * time.Second)
			}
			if err != nil {
				log.Printf("All retries failed for keyword '%s'\n", cleanKeyword)
				return
			}

			ProcessFileForKeyword(cleanKeyword, extensions, sessionCookie, results, concurrencyLimit)
		}

		wg.Wait()

	}

}

func readSessionCookie(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open session cookie file: %w", err)
	}
	defer file.Close()

	var sb strings.Builder
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		sb.WriteString(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading session cookie file: %w", err)
	}

	return strings.TrimSpace(sb.String()), nil
}

func SaveResults(results []map[string]interface{}, outputFile string) error {
	fmt.Println("These are the results", results)
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

	// fmt.Printf("Results saved to %s\n", outputFile)
	return nil
}

func ProcessFileForKeyword(keyword string, extensions map[string][]string, sessionCookie string, results []map[string]interface{}, concurrencyLimit int) {
	errorschan := make(chan error, len(files))
	semaphore := make(chan struct{}, concurrencyLimit)

	done := make(chan struct{})
	defer close(done)

	for _, fileInfo := range files {

		wg.Add(1)

		if fileInfo.Size > 50*1024*1024 {
			errorschan <- fmt.Errorf("skipping large file: %s", fileInfo.Filename)
			continue
		}

		processingColor := color.New(color.FgGreen).PrintlnFunc()
		go func(file api.FileInfo) {
			semaphore <- struct{}{}
			// processingColor("Starting a goroutine...")

			// defer func() {
			// 	processingColor("Exiting goroutine...")
			// }()
			defer wg.Done()

			start := time.Now()
			result := download.ProcessFile(file, extensions)
			processingColor("∟ File processed in →", time.Since(start), "\n")
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
	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Intn(900) + 100

	fileJSONName = fmt.Sprintf("results-%s-%d.json", keyword, randomNumber)
	fmt.Print("The proccess last less than 300 seconds. Saving current results...")
	mutex.Lock()
	err = SaveResults(results, fileJSONName)
	mutex.Unlock()
	if err != nil {
		log.Printf("Error saving results for keyword '%s': %v", keyword, err)
	}
	fmt.Println("Results saved to", fileJSONName)
	if err != nil {
		log.Printf("Error saving final results for keyword '%s': %v", keyword, err)
	}

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

}
