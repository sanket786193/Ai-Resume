package ocr

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ocrService implements the OCR Service interface
type ocrService struct {
	uploadDir string
}

// NewService creates a new OCR service instance
func NewService(uploadDir string) Service {
	return &ocrService{
		uploadDir: uploadDir,
	}
}

// ExtractText runs Tesseract OCR on an image file, PDF, or URL and returns the extracted text
func (s *ocrService) ExtractText(filePath, language string) (string, error) {
	if language == "" {
		language = "eng"
	}

	// Check if filePath is a URL (starts with http:// or https://)
	var localPath string
	if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
		// Download the file to a temporary file
		tempFile, err := s.downloadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to download file: %w", err)
		}
		defer os.Remove(tempFile) // Clean up temp file
		localPath = tempFile
	} else {
		localPath = filePath
	}

	// Check if it's a PDF
	if s.IsPDF(localPath) {
		return s.extractTextFromPDF(localPath, language)
	}

	// Process as image
	return s.extractTextFromImage(localPath, language)
}

// extractTextFromImage runs Tesseract OCR on an image file
func (s *ocrService) extractTextFromImage(imagePath, language string) (string, error) {
	// Tesseract command: tesseract <input> stdout -l <language>
	cmd := exec.Command("tesseract", imagePath, "stdout", "-l", language)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("tesseract error: %v, stderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// extractTextFromPDF converts PDF pages to images and extracts text from each page
func (s *ocrService) extractTextFromPDF(pdfPath, language string) (string, error) {
	// Check if pdftoppm is available
	if _, err := exec.LookPath("pdftoppm"); err != nil {
		return "", fmt.Errorf("pdftoppm not found. Please install poppler-utils: https://poppler.freedesktop.org/")
	}

	// Create temporary directory for PDF page images
	tempDir, err := os.MkdirTemp(s.uploadDir, "pdf-pages-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up temp directory

	// Convert PDF pages to images using pdftoppm
	// pdftoppm -png -r 300 input.pdf output_prefix
	outputPrefix := filepath.Join(tempDir, "page")
	cmd := exec.Command("pdftoppm", "-png", "-r", "300", pdfPath, outputPrefix)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pdftoppm error: %v, stderr: %s", err, stderr.String())
	}

	// Get all generated page images
	files, err := filepath.Glob(outputPrefix + "-*.png")
	if err != nil {
		return "", fmt.Errorf("failed to find PDF page images: %w", err)
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no pages found in PDF")
	}

	// Process each page and combine text
	var allText []string
	for i, pageFile := range files {
		pageText, err := s.extractTextFromImage(pageFile, language)
		if err != nil {
			return "", fmt.Errorf("failed to extract text from page %d: %w", i+1, err)
		}
		if pageText != "" {
			allText = append(allText, fmt.Sprintf("--- Page %d ---\n%s", i+1, pageText))
		}
	}

	return strings.Join(allText, "\n\n"), nil
}

// downloadFile downloads a file (image or PDF) from a URL to a temporary file
func (s *ocrService) downloadFile(url string) (string, error) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp(s.uploadDir, "ocr-*.tmp")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	// Copy the response body to the temporary file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

// ValidateImageFormat checks if the file extension is a valid image format
func (s *ocrService) ValidateImageFormat(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".gif":  true,
		".bmp":  true,
		".tiff": true,
		".webp": true,
	}
	return validExts[ext]
}

// ValidateFileFormat checks if the file extension is a valid format (images or PDF)
func (s *ocrService) ValidateFileFormat(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".gif":  true,
		".bmp":  true,
		".tiff": true,
		".webp": true,
		".pdf":  true,
	}
	return validExts[ext]
}

// IsPDF checks if the file is a PDF
func (s *ocrService) IsPDF(filename string) bool {
	return strings.ToLower(filepath.Ext(filename)) == ".pdf"
}

// EnsureUploadDir creates the upload directory if it doesn't exist
func (s *ocrService) EnsureUploadDir() error {
	return os.MkdirAll(s.uploadDir, 0755)
}

// GetImagePath returns the full path for an image file
func (s *ocrService) GetImagePath(filename string) string {
	return filepath.Join(s.uploadDir, filename)
}

// CheckTesseractAvailable checks if Tesseract is installed
func CheckTesseractAvailable() bool {
	_, err := exec.LookPath("tesseract")
	return err == nil
}

// CheckPopplerAvailable checks if pdftoppm (poppler-utils) is installed
func CheckPopplerAvailable() bool {
	_, err := exec.LookPath("pdftoppm")
	return err == nil
}
