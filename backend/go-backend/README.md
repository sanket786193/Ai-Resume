# Resume Backend - OCR Service

A Go-based backend service for OCR (Optical Character Recognition) processing with Cloudinary integration for file storage.

## Features

- **OCR Processing**: Extract text from images using Tesseract OCR
- **Cloudinary Integration**: Store and manage images in the cloud
- **RESTful API**: Clean API endpoints for image upload and processing
- **Flexible Storage**: Supports both Cloudinary and local file storage

## Project Structure

```
go-backend/
├── api/
│   ├── handlers/          # HTTP handlers
│   │   └── health.go
│   └── routes.go          # Route definitions
├── internal/
│   ├── config/            # Configuration management
│   │   └── config.go
│   └── service/
│       └── ocr/           # OCR service layer
│           ├── handler.go        # OCR HTTP handlers
│           ├── service.go        # OCR business logic
│           ├── storage.go        # Local file storage
│           ├── cloudinary_storage.go  # Cloudinary storage
│           └── model.go          # Data models
├── uploads/               # Local upload directory (if not using Cloudinary)
├── main.go               # Application entry point
└── go.mod
```

## Setup

### Prerequisites

1. **Go 1.24+** installed
2. **Tesseract OCR** installed:
   - Windows: Download from [GitHub](https://github.com/UB-Mannheim/tesseract/wiki)
   - macOS: `brew install tesseract`
   - Linux: `sudo apt-get install tesseract-ocr`

### Installation

1. Clone the repository and navigate to the backend:
```bash
cd go-backend
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables (create a `.env` file or export them):

**For Cloudinary (Recommended):**
```bash
export CLOUDINARY_CLOUD_NAME=your_cloud_name
export CLOUDINARY_API_KEY=your_api_key
export CLOUDINARY_API_SECRET=your_api_secret
export CLOUDINARY_FOLDER=resume-ocr
export CLOUDINARY_ENABLED=true
```

**For Local Storage:**
```bash
export CLOUDINARY_ENABLED=false
export UPLOAD_DIR=uploads
```

**Server Configuration:**
```bash
export PORT=8080
```

## API Endpoints

### Health Check
```
GET /health
```
Returns server status.

**Response:**
```json
{
  "status": "UP",
  "service": "resume-backend"
}
```

### Upload and Process Image
```
POST /ocr/image
Content-Type: multipart/form-data
```

**Form Data:**
- `image` (file, required): Image file to process
- `language` (string, optional): OCR language code (default: "eng")

**Response:**
```json
{
  "text": "Extracted text from image",
  "filename": "image.png",
  "language": "eng",
  "path": "https://res.cloudinary.com/..."
}
```

### Test OCR (Process existing file)
```
GET /ocr/test/:filename?language=eng
```

**Parameters:**
- `filename` (path parameter): Name of the file in storage
- `language` (query parameter, optional): OCR language code

**Response:**
```json
{
  "text": "Extracted text from image",
  "filename": "test1.png",
  "language": "eng",
  "path": "https://res.cloudinary.com/..."
}
```

## Usage Examples

### Using cURL

**Upload and process an image:**
```bash
curl -X POST http://localhost:8080/ocr/image \
  -F "image=@/path/to/image.png" \
  -F "language=eng"
```

**Process existing file:**
```bash
curl http://localhost:8080/ocr/test/test1.png?language=eng
```

### Using JavaScript (Fetch API)

```javascript
const formData = new FormData();
formData.append('image', fileInput.files[0]);
formData.append('language', 'eng');

fetch('http://localhost:8080/ocr/image', {
  method: 'POST',
  body: formData
})
.then(response => response.json())
.then(data => console.log(data));
```

## Architecture

### Service Layer
- **OCR Service**: Handles OCR processing logic, supports both local files and URLs
- **Storage Interface**: Abstract storage layer supporting multiple implementations
  - **Local Storage**: File-based storage on disk
  - **Cloudinary Storage**: Cloud-based image storage and CDN

### Handler Layer
- **OCR Handler**: HTTP handlers for OCR endpoints
- **Health Handler**: Health check endpoint

### Configuration
- Environment-based configuration
- Supports both Cloudinary and local storage
- Configurable upload directory and Cloudinary folder

## Development

### Running the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080` (or the port specified in `PORT` environment variable).

### Building

```bash
go build -o resume-backend
./resume-backend
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `UPLOAD_DIR` | Local upload directory | `uploads` |
| `CLOUDINARY_ENABLED` | Enable Cloudinary storage | `true` |
| `CLOUDINARY_CLOUD_NAME` | Cloudinary cloud name | - |
| `CLOUDINARY_API_KEY` | Cloudinary API key | - |
| `CLOUDINARY_API_SECRET` | Cloudinary API secret | - |
| `CLOUDINARY_FOLDER` | Cloudinary folder path | `resume-ocr` |

## Error Handling

The API returns standard HTTP status codes:
- `200 OK`: Success
- `400 Bad Request`: Invalid request (missing file, invalid format)
- `404 Not Found`: File not found
- `500 Internal Server Error`: Server error (OCR failure, storage error)

## License

MIT
