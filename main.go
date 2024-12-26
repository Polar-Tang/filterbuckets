package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/Polar-Tang/filterbuckets/processing"
)

// ToolConfig holds configuration for the tool
type ToolConfig struct {
	KeywordsFile   string
	ExtensionsFile string
}

// Main entry point
func main() {
	config := parseFlags()
	keywords := readLinesFromFile(config.KeywordsFile)
	extensions, err := readExtensionsMap(config.ExtensionsFile)
	if err != nil {
		fmt.Printf("Error reading extensions file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded %d keywords.\n", len(keywords))
	fmt.Printf("Loaded %d extensions.\n", len(extensions))

	// Placeholder: Call your processing function here
	processing.ProcessFiles(keywords, extensions)
}

// parseFlags parses and validates command-line arguments
func parseFlags() ToolConfig {
	keywordsFile := flag.String("w", "", "Path to the file containing keywords.")
	extensionsFile := flag.String("x", "", "Path to the file containing extensions.")
	flag.Parse()

	if *extensionsFile == "" {
		fmt.Println("Error: --extensions flag is required.")
		flag.Usage()
		os.Exit(1)
	}

	return ToolConfig{
		KeywordsFile:   *keywordsFile,
		ExtensionsFile: *extensionsFile,
	}
}

func readExtensionsMap(filePath string) (map[string][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	extensionsMap := make(map[string][]string)

	if err := decoder.Decode(&extensionsMap); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	return extensionsMap, nil
}

// readLinesFromFile reads lines from a given file and returns them as a slice of strings
func readLinesFromFile(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		os.Exit(1)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		os.Exit(1)
	}
	return lines
}

// processFiles is a placeholder for your main processing lo
