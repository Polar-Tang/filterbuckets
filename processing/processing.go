package processing

// IMPORTS
import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Polar-Tang/filterbuckets/api"
	"github.com/Polar-Tang/filterbuckets/download"

	"github.com/fatih/color"
)

func ProcessFiles(keywords []string, extensions map[string][]string) {
	sessionCookie := "54e7fe8c2aa1dd504b9be39fa3466f10"
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
		cleankeyword := strings.TrimSpace(keyword)
		cleankeyword = strings.Trim(cleankeyword, `"`)
		cleankeyword = strings.TrimRight(cleankeyword, `",`)

		outputFile, err := createOutputFile(cleankeyword)
		if err != nil {
			fmt.Printf("Failed to create output file: %v\n", err)
			continue
		}
		fmt.Printf("Searching for files with keyword: %s\n", cleankeyword)

		var files []api.FileInfo
		maxRetries := 3
		for retries := 0; retries < maxRetries; retries++ {
			files, err = api.QueryFiles(sessionCookie, []string{cleankeyword}, extensions)
			if err == nil {
				break
			}
			log.Printf("Retry %d/%d for keyword '%s' failed: %v", retries+1, maxRetries, cleankeyword, err)
			time.Sleep(2 * time.Second)
		}
		if err != nil {
			log.Printf("All retries failed for keyword '%s'\n", cleankeyword)
			continue
		}

		concurrencyLimit := 6
		semaphore := make(chan struct{}, concurrencyLimit)
		errorschan := make(chan error, len(files))
		results := make([]map[string]interface{}, 0)
		var wg sync.WaitGroup
		var mutex sync.Mutex

		ticker := time.NewTicker(300 * time.Second)
		defer ticker.Stop()

		// 🩵🩵🩵🩵🩵🩵🩵🩵🩵🩵
		tickerColor := color.New(color.FgBlue).PrintlnFunc()
		go func() {

			buf := make([]byte, 1<<16) // Large enough buffer for goroutines dump
			runtime.Stack(buf, true)
			fmt.Printf("%s\n", buf)

			tickerColor("Starting a goroutine...")
			defer tickerColor("Exiting goroutine...")
			for range ticker.C {
				tickerColor("Locking mutex...\n")
				mutex.Lock()
				err := SaveResults(results, outputFile)
				if err != nil {
					errorschan <- fmt.Errorf("error saving periodic results for keyword '%s': %v", keyword, err)
				}
				tickerColor("Unocking mutex...\n")
				mutex.Unlock()
			}
		}()

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
		wg.Wait()
		mutex.Lock()
		err = SaveResults(results, outputFile)
		if err != nil {
			log.Printf("Error saving final results for keyword '%s': %v", keyword, err)
		}
		mutex.Unlock()
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
