package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"pdf_greyhat_go/api"
	"pdf_greyhat_go/processing"

	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api" // Alias for pdfcpu API
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

	// res_content, err := readFileContent(tmpFile.Name())

	// if err != nil {
	// 	fmt.Printf("Failed to read extracted content from %s without use pdfcpu: %v\n", tmpFile.Name(), err)
	// 	return nil
	// }
	// fmt.Println("Extracted Content:\n", res_content) // Print the extracted content
	// ---------------------------------------------------------------------------------------------

	// Step 4: Extract content from the PDF

	// create the directory
	outputDir, err := os.MkdirTemp("", "pdf_extracted_*")

	if err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		return nil
	}
	defer os.Remove(outputDir) // Clean up the extracted text file

	// extract the content from the directory

	err = pdfcpuapi.ExtractContentFile(tmpFile.Name(), outputDir, nil, nil)
	if err != nil {
		fmt.Printf("PDFCPU extraction failed for file %s: %v\n", tmpFile.Name(), err)

		// Assume failure is due to image content, process with OCR
		fmt.Println("Falling back to OCR for image-based PDF...")
		ocrOutput := filepath.Join(outputDir, "ocr_output")
		err := processing.RunTesseract(tmpFile.Name(), ocrOutput)
		if err != nil {
			fmt.Printf("OCR failed for %s: %v\n", tmpFile.Name(), err)
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
	//fmt.Println("Combined content from extracted text files:", contentpdf) // debuging

	// Opens the tmp file and reads it line by line using bufio.Scanner.

	// // Set the path to the PDF and read images
	// client.SetImage(tmpFile.Name())
	// text, err := client.Text()
	// if err != nil {
	// 	fmt.Printf("Failed to extract text from PDF %s: %v\n", file.Filename, err)
	// 	return nil
	// }

	// ---------------------------------------------------------------------------------------------

	// Step 6: Analyze text for keywords
	// Accumulator
	keywordCounts := make(map[string]int)

	// iterates over the whole file, looking for our pdfKeywords (argument)
	for _, keyword := range pdfKeywords {
		// search for the keyword, in a insasitive case way, on the textContent, which is the text extracted from the output
		count := strings.Count(strings.ToLower(contentpdf), strings.ToLower(keyword)) // strings.Count is a built-in function from string package, used to count the words
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
// func readFileContent(filePath string) (string, error) {
// 	// Open the file
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer file.Close()

// 	// Use bufio to read the file line by line
// 	var sb strings.Builder
// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		sb.WriteString(scanner.Text() + "\n")
// 	}
// 	if err := scanner.Err(); err != nil {
// 		return "", err
// 	}

// 	return sb.String(), nil
// }

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
