package ocr

// OCRRequest represents the request for OCR processing
type OCRRequest struct {
	ImagePath string
	Language  string
}

// OCRResponse represents the response from OCR processing
type OCRResponse struct {
	Text     string `json:"text"`
	Filename string `json:"filename"`
	Language string `json:"language"`
	Path     string `json:"path,omitempty"`
}

// Service defines the OCR service interface
type Service interface {
	ExtractText(filePath, language string) (string, error)
	ValidateImageFormat(filename string) bool
	ValidateFileFormat(filename string) bool
	IsPDF(filename string) bool
}
