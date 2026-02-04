package ocr

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// Handler handles OCR-related HTTP requests
type Handler struct {
	service Service
	storage Storage
}

// NewHandler creates a new OCR handler
func NewHandler(service Service, storage Storage) *Handler {
	return &Handler{
		service: service,
		storage: storage,
	}
}

// ImageOCR handles POST /ocr/image - processes uploaded image or PDF
func (h *Handler) ImageOCR(c echo.Context) error {
	// Get uploaded file (can be image or PDF)
	file, err := c.FormFile("image")
	if err != nil {
		// Try alternative field name "file"
		file, err = c.FormFile("file")
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "File required. Use form field 'image' or 'file'",
			})
		}
	}

	// Validate file extension (images or PDF)
	if !h.service.ValidateFileFormat(file.Filename) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid file format. Supported: PNG, JPG, JPEG, GIF, BMP, TIFF, WEBP, PDF",
		})
	}

	// Check if PDF and verify poppler is available
	if h.service.IsPDF(file.Filename) {
		if !CheckPopplerAvailable() {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "PDF processing requires pdftoppm. Please install poppler-utils: https://poppler.freedesktop.org/",
			})
		}
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to open uploaded file",
		})
	}
	defer src.Close()

	// Save file (returns Cloudinary URL or local path)
	filePath, err := h.storage.SaveFile(src, file.Filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to save file: %v", err),
		})
	}

	// Get language parameter (default: eng)
	language := c.FormValue("language")
	if language == "" {
		language = "eng"
	}

	// Perform OCR on the file (service handles both URLs and local paths, images and PDFs)
	text, err := h.service.ExtractText(filePath, language)
	if err != nil {
		// Clean up uploaded file (only if it's a Cloudinary URL)
		if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
			h.storage.DeleteFile(filePath)
		}

		// Check if tesseract is available
		if !CheckTesseractAvailable() {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Tesseract OCR not found. Please install Tesseract: https://github.com/tesseract-ocr/tesseract",
			})
		}

		// Check if poppler is needed and available
		if h.service.IsPDF(file.Filename) && !CheckPopplerAvailable() {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "PDF processing requires pdftoppm. Please install poppler-utils: https://poppler.freedesktop.org/",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("OCR failed: %v", err),
		})
	}

	// Clean up uploaded file (only if it's a Cloudinary URL)
	if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
		h.storage.DeleteFile(filePath)
	}

	return c.JSON(http.StatusOK, OCRResponse{
		Text:     text,
		Filename: file.Filename,
		Language: language,
		Path:     filePath, // Include Cloudinary URL in response
	})
}

// TestOCR handles GET /ocr/test/:filename - processes image or PDF from uploads folder
func (h *Handler) TestOCR(c echo.Context) error {
	filename := c.Param("filename")
	if filename == "" {
		filename = "test1.png" // Default to test1.png
	}

	// Get language parameter (default: eng)
	language := c.QueryParam("language")
	if language == "" {
		language = "eng"
	}

	// Get file path using storage
	filePath := h.storage.GetImagePath(filename)

	// For Cloudinary, we need to construct the full URL
	// For local storage, check if file exists
	if !strings.HasPrefix(filePath, "http://") && !strings.HasPrefix(filePath, "https://") {
		if !h.storage.FileExists(filePath) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": fmt.Sprintf("File not found: %s", filename),
			})
		}
	} else {
		// For Cloudinary URLs, check if file exists
		if !h.storage.FileExists(filePath) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": fmt.Sprintf("File not found: %s", filename),
			})
		}
	}

	// Check if PDF and verify poppler is available
	if h.service.IsPDF(filename) {
		if !CheckPopplerAvailable() {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "PDF processing requires pdftoppm. Please install poppler-utils: https://poppler.freedesktop.org/",
			})
		}
	}

	// Perform OCR on the file (service handles both URLs and local paths, images and PDFs)
	text, err := h.service.ExtractText(filePath, language)
	if err != nil {
		// Check if tesseract is available
		if !CheckTesseractAvailable() {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Tesseract OCR not found. Please install Tesseract: https://github.com/tesseract-ocr/tesseract",
			})
		}

		// Check if poppler is needed and available
		if h.service.IsPDF(filename) && !CheckPopplerAvailable() {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "PDF processing requires pdftoppm. Please install poppler-utils: https://poppler.freedesktop.org/",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("OCR failed: %v", err),
		})
	}

	return c.JSON(http.StatusOK, OCRResponse{
		Text:     text,
		Filename: filename,
		Language: language,
		Path:     filePath,
	})
}
