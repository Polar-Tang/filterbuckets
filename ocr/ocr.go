package ocr

import (
	"fmt"
	"os/exec"
)

// ExtractText processes an image file with Tesseract OCR and saves the output as a text file
func RunTesseract(imagePath string, outputTextPath string) error {
	fmt.Printf("Running Tesseract OCR on %s, outputting to %s.txt\n", imagePath, outputTextPath)

	// Prepare the Tesseract command
	cmd := exec.Command("tesseract", imagePath, outputTextPath)
	cmd.Stdout = nil // Silence stdout
	cmd.Stderr = nil // Silence stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("tesseract OCR failed: %w", err)
	}

	fmt.Printf("OCR completed successfully. Text saved to %s.txt\n", outputTextPath)
	return nil
}
