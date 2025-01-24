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
	"github.com/fatih/color"

	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api"
)

func ProcessFile(file api.FileInfo, extensionKeywords map[string][]string) map[string]interface{} {

	bucketTransport := &http.Transport{
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 5,
		IdleConnTimeout:     5 * time.Second,
		DisableKeepAlives:   true,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	bucketClient := &http.Client{
		Transport: bucketTransport,
		Timeout:   15 * time.Second,
	}

	response, err := bucketClient.Get(file.URL)

	if err != nil {
		fmt.Printf("Failed to download file %s: %v\n", file.URL, err)
		return nil
	}

	defer response.Body.Close()
	fmt.Printf("Processing url %s\n", file.URL)

	tmpFile, err := os.CreateTemp("", "*"+filepath.Ext(file.Filename))

	if err != nil {
		fmt.Printf("Failed to create temp file for %s: %v\n", file.URL, err)
		return nil
	}

	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, response.Body)
	if err != nil {
		fmt.Printf("Failed to save file %s: %v\n", file.URL, err)
		return nil
	}

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

func processPlainText(filePath string, keywords []string, file api.FileInfo) map[string]interface{} {
	// fmt.Println("Processing ", file.FullPath)
	content, err := readFileContent(filePath)
	if err != nil {
		fmt.Printf("Failed to read content from %s: %v\n", file.Filename, err)
		return nil
	}

	processingColor := color.New(color.FgRed).PrintlnFunc()
	processingColorGreen := color.New(color.FgGreen).PrintlnFunc()

	keywordCounts := countKeywords(content, keywords)

	if keywordCounts != nil {
		processingColorGreen("∟ Keywords found")
	} else {
		processingColor("∟ No keywords found")
		return nil
	}
	// fmt.Print(keywordCounts)

	return map[string]interface{}{
		"url":      file.URL,
		"filename": file.Filename,
		"keywords": keywordCounts,
	}
}

func readFileContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var sb strings.Builder
	scanner := bufio.NewScanner(file)

	bufferSize := 10 * 1024 * 1024
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

func countKeywords(content string, keywords []string) map[string]int {

	keywordCounts := make(map[string]int)

	content = strings.ToLower(content)

	for _, keyword := range keywords {
		pattern := "(?i)" + keyword

		re, err := regexp.Compile(pattern)
		if err != nil {
			fmt.Printf("Failed to compile regex for keyword '%s': %v\n", keyword, err)
			continue
		}

		matches := re.FindAllStringIndex(content, -1)
		if len(matches) > 0 {

			keywordCounts[keyword] = len(matches)
		}
	}

	if len(keywordCounts) == 0 {
		return nil
	}

	return keywordCounts
}

func processPDF(filePath string, keywords []string, file api.FileInfo) map[string]interface{} {

	// fmt.Println("Processing: ", file.FullPath)

	outputDir, err := os.MkdirTemp("", "pdf_extracted_*")

	if err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		return nil
	}
	defer os.RemoveAll(outputDir)

	err = pdfcpuapi.ExtractContentFile(filePath, outputDir, nil, nil)
	if err != nil {
		fmt.Printf("PDFCPU extraction failed for file %s: %v\n", outputDir, err)

		fmt.Println("Falling back to OCR for image-based PDF...")

		return nil
	}

	entries, err := os.ReadDir(outputDir)
	if err != nil {
		fmt.Printf("Failed to list directory %s: %v\n", outputDir, err)
		return nil
	}

	for _, entry := range entries {
		fmt.Println("-", entry.Name())
	}

	contentpdf, err := readExtractedText(outputDir)
	if err != nil {
		fmt.Printf("Failed to read extracted text from directory %s: %v\n", outputDir, err)
		return nil
	}
	keywordCounts := countKeywords(contentpdf, keywords)

	processingColor := color.New(color.FgRed).PrintlnFunc()
	processingColorGreen := color.New(color.FgGreen).PrintlnFunc()

	if keywordCounts != nil {
		processingColorGreen("∟ Keywords found")
	} else {
		processingColor("∟ No keywords found")
		return nil
	}
	fmt.Print(keywordCounts)
	return map[string]interface{}{
		"url":      file.URL,
		"filename": file.Filename,
		"keywords": keywordCounts,
	}

}

func readExtractedText(dir string) (string, error) {
	var sb strings.Builder

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			path := filepath.Join(dir, entry.Name())

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

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return false, err
	}

	mimeType := http.DetectContentType(buffer)
	return strings.HasPrefix(mimeType, "image/"), nil
}
