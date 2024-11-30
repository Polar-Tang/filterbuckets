package download

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"pdf_processing_go/api" // Import the api package for FileInfo struct

	"github.com/otiai10/gosseract/v2"
	// Alias for pdfcpu API
)

// ProcessFile downloads and analyzes the PDF file for keywords
func ProcessFile(file api.FileInfo, pdfKeywords []string) map[string]interface{} {
	// fmt.Println("Processing file:", file.URL)

	// ---------------------------------------------------------------------------------------------
	// Step 1: Download the PDF file
	response, err := http.Get(file.URL) // download the url with a simple get
	// A simple error handling
	if err != nil {
		// fmt.Printf("Failed to download file %s: %v\n", file.URL, err)
		return nil
	}
	// Close the response body
	defer response.Body.Close()

	// ---------------------------------------------------------------------------------------------

	// 2) Create a temporary file, referenced by name
	tmpFile, err := os.CreateTemp("", "*.pdf")
	// A simple error handling
	if err != nil {
		// fmt.Printf("Failed to create temp file for %s: %v\n", file.URL, err)
		return nil
	}
	// remove the TMP after the execution
	defer os.Remove(tmpFile.Name()) //  using Name, builtin function from the os.file package, this returns the full path to the temporary file created

	// ---------------------------------------------------------------------------------------------

	// 3) copy th TMP
	// Ignoring the return value with the blank identifier, just copy the TMP using io.Copy
	_, err = io.Copy(tmpFile, response.Body) // copies the response body to the tmp file
	if err != nil {
		fmt.Printf("Failed to save file %s: %v\n", file.URL, err)
		return nil
	}

	// ---------------------------------------------------------------------------------------------

	// Step 4: Extract content from the PDF
	client := gosseract.NewClient()
	defer client.Close()

	// Set the path to the PDF and read images
	client.SetImage(tmpFile.Name())
	text, err := client.Text()
	if err != nil {
		fmt.Printf("Failed to extract text from PDF %s: %v\n", file.Filename, err)
		return nil
	}

	// ---------------------------------------------------------------------------------------------

	// Step 6: Analyze text for keywords
	// Accumulator
	keywordCounts := make(map[string]int)

	// iterates over the whole file, looking for our pdfKeywords (argument)
	for _, keyword := range pdfKeywords {
		// search for the keyword, in a insasitive case way, on the textContent, which is the text extracted from the output
		count := strings.Count(strings.ToLower(text), strings.ToLower(keyword)) // strings.Count is a built-in function from string package, used to count the words
		// save them in the accumulator
		keywordCounts[keyword] = count
	}

	// ---------------------------------------------------------------------------------------------
	fmt.Printf("Keyword counts for %s: %+v\n", file.Filename, keywordCounts)

	// Step 7: Return results only if keywords are found
	anyKeywordFound := false
	for _, count := range keywordCounts {
		if count > 0 {
			anyKeywordFound = true
			break
		}
	}

	if !anyKeywordFound {
		fmt.Printf("No keywords found in %s\n", file.Filename)
		return nil // Skip files with no keyword matches
	}

	// Step 7: Return results, with the keword counts
	return map[string]interface{}{
		"url":      file.URL,
		"filename": file.Filename,
		"keywords": keywordCounts,
	}
}

// The function to read the file in the step 3
func readFileContent(filePath string) (string, error) {
	// opens the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// uses bufio to read the file line by line
	var sb strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sb.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	// returns the content
	return sb.String(), nil
}
