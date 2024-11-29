package download

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api" // Alias for pdfcpu API
	"pdf_processing_go/api"                      // Import the api package for FileInfo struct
)

// ProcessFile downloads and analyzes the PDF file for keywords
func ProcessFile(file api.FileInfo, pdfKeywords []string) map[string]interface{} {
	fmt.Println("Processing file:", file.URL)

	// Step 1: Download the PDF file
	response, err := http.Get(file.URL)
	if err != nil {
		fmt.Printf("Failed to download file %s: %v\n", file.URL, err)
		return nil
	}
	defer response.Body.Close()

	// Save file to a temporary local file
	tmpFile, err := os.CreateTemp("", "*.pdf")
	if err != nil {
		fmt.Printf("Failed to create temp file for %s: %v\n", file.URL, err)
		return nil
	}
	defer os.Remove(tmpFile.Name()) // Clean up the temp file after processing

	_, err = io.Copy(tmpFile, response.Body)
	if err != nil {
		fmt.Printf("Failed to save file %s: %v\n", file.URL, err)
		return nil
	}

	// Step 2: Extract content from the PDF
	outputFile := tmpFile.Name() + "_extracted.txt"
	err = pdfcpuapi.ExtractContentFile(tmpFile.Name(), outputFile, nil, nil)
	if err != nil {
		fmt.Printf("Failed to extract content from PDF %s: %v\n", file.Filename, err)
		return nil
	}
	defer os.Remove(outputFile) // Clean up the extracted text file

	// Step 3: Read extracted text
	textContent, err := readFileContent(outputFile)
	if err != nil {
		fmt.Printf("Failed to read extracted content from %s: %v\n", outputFile, err)
		return nil
	}

	// Step 4: Analyze text for keywords
	keywordCounts := make(map[string]int)
	for _, keyword := range pdfKeywords {
		count := strings.Count(strings.ToLower(textContent), strings.ToLower(keyword))
		keywordCounts[keyword] = count
	}

	// Step 5: Return results
	return map[string]interface{}{
		"url":      file.URL,
		"filename": file.Filename,
		"keywords": keywordCounts,
	}
}

// Helper function to read file content into a string
func readFileContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var sb strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sb.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return sb.String(), nil
}
