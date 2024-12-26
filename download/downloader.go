package download

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Polar-Tang/filterbuckets/api"
	"github.com/Polar-Tang/filterbuckets/ocr"

	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api" // Alias for pdfcpu API
)

// ProcessFile downloads and analyzes the PDF file for keywords
func ProcessFile(file api.FileInfo, extensionKeywords map[string][]string) map[string]interface{} {

	// fmt.Println("Processing file:", file.URL)
	// transport for the buckets
	bucketTransport := &http.Transport{
		MaxIdleConns:        50,                                    // Adjust as per workload
		MaxIdleConnsPerHost: 5,                                     // Limit per host
		IdleConnTimeout:     5 * time.Second,                       // Free idle connections quickly
		DisableKeepAlives:   true,                                  // Avoid reusing connections
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true}, // Skip SSL verification
	}

	// Create an HTTP client for bucket queries
	bucketClient := &http.Client{
		Transport: bucketTransport,
		Timeout:   15 * time.Second, // Overall timeout for the bucket queries
	}

	// ---------------------------------------------------------------------------------------------
	// Step 1: Download the PDF file
	response, err := bucketClient.Get(file.URL) // download the url with a simple get
	// A simple error handling
	if err != nil {
		fmt.Printf("Failed to download file %s: %v\n", file.URL, err)
		return nil
	}
	// Close the response body
	defer response.Body.Close()

	// ---------------------------------------------------------------------------------------------

	// 2) Create a temporary file, referenced by name
	tmpFile, err := os.CreateTemp("", "*"+filepath.Ext(file.Filename))
	// A simple error handling
	if err != nil {
		fmt.Printf("Failed to create temp file for %s: %v\n", file.URL, err)
		return nil
	}

	// remove the TMP after the execution
	defer os.Remove(tmpFile.Name()) //  using Name, builtin function from the os.file package, this returns the full path to the temporary file created

	// ---------------------------------------------------------------------------------------------

	// 3) copy the download to a TMP
	// Ignoring the return value with the blank identifier, just copy the TMP using io.Copy
	_, err = io.Copy(tmpFile, response.Body) // copies the response body to the tmp file
	if err != nil {
		fmt.Printf("Failed to save file %s: %v\n", file.URL, err)
		return nil
	}

	// 4) Determine the file extension

	extension := strings.TrimPrefix(filepath.Ext(tmpFile.Name()), ".")
	keywords, found := extensionKeywords[extension]

	if !found {
		fmt.Printf("Skipping unsupported file type: %s\n", extension)
		return nil
	}

	if extension == "pdf" {
		return processPDF(tmpFile.Name(), keywords, file)
	} else {
		return processPlainText(tmpFile.Name(), keywords, file)
	}

}

// --------------------------------------------------------------------------------------------

func processPlainText(filePath string, keywords []string, file api.FileInfo) map[string]interface{} {
	fmt.Println("Processing plain from the URL:", file.FullPath)
	content, err := readFileContent(filePath)
	if err != nil {
		fmt.Printf("Failed to read content from %s: %v\n", file.Filename, err)
		return nil
	}

	keywordCounts := countKeywords(content, keywords)

	if keywordCounts == nil || len(keywordCounts) == 0 {
		fmt.Printf("No keywords found in file: %s\n", file.Filename)
		return nil
	}

	return map[string]interface{}{
		"url":      file.URL,
		"filename": file.Filename,
		"keywords": keywordCounts,
	}
}

// --------------------------------------------------------------------------------------------

func readFileContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var sb strings.Builder
	scanner := bufio.NewScanner(file)

	// Increase the buffer size to handle large tokens
	bufferSize := 10 * 1024 * 1024 // 10 MB
	scanner.Buffer(make([]byte, bufferSize), bufferSize)

	for scanner.Scan() {
		sb.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		content, err := io.ReadAll(file)
		if err != nil {
			return "Something went wrong in io reader", err
		}
		return string(content), nil
	}
	return sb.String(), nil
}

// ---------------------------------------------------------------------------------------------

func countKeywords(content string, keywords []string) map[string]int {
	// Accumulator
	keywordCounts := make(map[string]int)
	// iterates over the whole file, looking for our pdfKeywords (argument)
	content = strings.ToLower(content)

	for _, keyword := range keywords {
		// search for the keyword, in a insasitive case way, on the textContent, which is the text extracted from the output
		pattern := "(?i)" + keyword // "(?i)" makes it case-insensitive

		re, err := regexp.Compile(pattern)
		if err != nil {
			fmt.Printf("Failed to compile regex for keyword '%s': %v\n", keyword, err)
			continue
		}

		matches := re.FindAllStringIndex(content, -1)
		if len(matches) > 0 {
			// Save the count of matches
			keywordCounts[keyword] = len(matches)
		}
	}

	if len(keywordCounts) == 0 {
		return nil
	}

	return keywordCounts
}

// ---------------------------------------------------------------------------------------------

func processPDF(filePath string, keywords []string, file api.FileInfo) map[string]interface{} {
	// create the directory
	fmt.Println("Processing plain from the URL:", file.FullPath)

	outputDir, err := os.MkdirTemp("", "pdf_extracted_*")

	if err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		return nil
	}
	defer os.RemoveAll(outputDir) // Clean up the extracted text file

	// extract the content from the directory
	err = pdfcpuapi.ExtractContentFile(filePath, outputDir, nil, nil)
	if err != nil {
		fmt.Printf("PDFCPU extraction failed for file %s: %v\n", outputDir, err)

		// Assume failure is due to image content, process with OCR
		fmt.Println("Falling back to OCR for image-based PDF...")
		ocrOutput := filepath.Join(outputDir, "ocr_output")
		err := ocr.RunTesseract(outputDir, ocrOutput)
		if err != nil {
			fmt.Printf("OCR failed for %s: %v\n", outputDir, err)
			return nil
		}
		fmt.Printf("OCR output saved to: %s.txt\n", ocrOutput)
		return nil
	}

	// Step 5: Read extracted text

	entries, err := os.ReadDir(outputDir)
	if err != nil {
		fmt.Printf("Failed to list directory %s: %v\n", outputDir, err)
		return nil
	}
	fmt.Println("Files in output directory:")
	for _, entry := range entries {
		fmt.Println("-", entry.Name())
	}

	// Read each page in the content
	contentpdf, err := readExtractedText(outputDir)
	if err != nil {
		fmt.Printf("Failed to read extracted text from directory %s: %v\n", outputDir, err)
		return nil
	}
	keywordCounts := countKeywords(contentpdf, keywords)

	if keywordCounts == nil || len(keywordCounts) == 0 {
		fmt.Printf("No keywords found in file: %s\n", file.Filename)
		return nil
	}
	return map[string]interface{}{
		"url":      file.URL,
		"filename": file.Filename,
		"keywords": keywordCounts,
	}

}

// ---------------------------------------------------------------------------------------------

func readExtractedText(dir string) (string, error) {
	var sb strings.Builder

	// List all entries in the directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	// Iterate over each file in the directory
	for _, entry := range entries {
		if !entry.IsDir() { // Skip subdirectories
			path := filepath.Join(dir, entry.Name())

			// Check if the file is an image
			isImage, err := isImageMime(path)
			if err != nil {
				fmt.Printf("Failed to check MIME type for %s: %v\n", entry.Name(), err)
				continue
			}
			if isImage {
				fmt.Println("Found an image file:", entry.Name())
				continue
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return "", fmt.Errorf("failed to read file %s: %w", path, err)
			}
			sb.Write(content)
		}
	}

	return sb.String(), nil
}

func isImageMime(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read file header to detect content type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return false, err
	}

	// Detect the MIME type
	mimeType := http.DetectContentType(buffer)
	return strings.HasPrefix(mimeType, "image/"), nil
}
