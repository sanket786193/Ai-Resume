package ocr

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// cloudinaryStorage implements Storage interface using Cloudinary
type cloudinaryStorage struct {
	cld       *cloudinary.Cloudinary
	ctx       context.Context
	folder    string
	cloudName string
}

// NewCloudinaryStorage creates a new Cloudinary storage instance
func NewCloudinaryStorage(cloudName, apiKey, apiSecret, folder string) (Storage, error) {
	ctx := context.Background()
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}

	return &cloudinaryStorage{
		cld:       cld,
		ctx:       ctx,
		folder:    folder,
		cloudName: cloudName,
	}, nil
}

// SaveFile uploads a file to Cloudinary and returns the public URL
func (s *cloudinaryStorage) SaveFile(src io.Reader, filename string) (string, error) {
	// Generate unique public ID from filename
	publicID := s.getPublicID(filename)

	overwrite := true

	// Upload to Cloudinary
	uploadParams := uploader.UploadParams{
		PublicID:     publicID,
		Folder:       s.folder,
		ResourceType: "image",
		Overwrite:    &overwrite,
	}

	result, err := s.cld.Upload.Upload(s.ctx, src, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}

	// Return the secure URL
	return result.SecureURL, nil
}

// DeleteFile deletes a file from Cloudinary
func (s *cloudinaryStorage) DeleteFile(filepath string) error {
	// Extract public ID from URL or use the path directly
	publicID := s.extractPublicID(filepath)

	// Delete from Cloudinary
	deleteParams := uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image",
	}

	_, err := s.cld.Upload.Destroy(s.ctx, deleteParams)
	if err != nil {
		return fmt.Errorf("failed to delete from Cloudinary: %w", err)
	}

	return nil
}

// FileExists checks if a file exists in Cloudinary (by attempting to get its info)
func (s *cloudinaryStorage) FileExists(filepath string) bool {
	publicID := s.extractPublicID(filepath)

	// Try to get resource info
	_, err := s.cld.Admin.Asset(s.ctx, admin.AssetParams{
		PublicID: publicID,
	})

	return err == nil
}

// GetImagePath returns the Cloudinary URL for a given filename
func (s *cloudinaryStorage) GetImagePath(filename string) string {
	publicID := s.getPublicID(filename)
	fullPublicID := publicID
	if s.folder != "" {
		fullPublicID = fmt.Sprintf("%s/%s", s.folder, publicID)
	}

	// Construct Cloudinary URL
	// Format: https://res.cloudinary.com/{cloud_name}/image/upload/{public_id}
	return fmt.Sprintf("https://res.cloudinary.com/%s/image/upload/%s", s.cloudName, fullPublicID)
}

// getPublicID generates a public ID from filename (removes extension)
func (s *cloudinaryStorage) getPublicID(filename string) string {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	// Replace spaces and special chars with underscores
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

// extractPublicID extracts public ID from Cloudinary URL or returns the path as-is
func (s *cloudinaryStorage) extractPublicID(path string) string {
	// If it's a Cloudinary URL, extract the public ID
	if strings.Contains(path, "cloudinary.com") {
		// Extract public ID from URL
		// Format: https://res.cloudinary.com/{cloud_name}/image/upload/{folder}/{public_id}.{ext}
		parts := strings.Split(path, "/")
		for i, part := range parts {
			if part == "upload" && i+1 < len(parts) {
				// Get everything after "upload"
				publicID := strings.Join(parts[i+1:], "/")
				// Remove extension
				ext := filepath.Ext(publicID)
				return strings.TrimSuffix(publicID, ext)
			}
		}
	}

	// If it's just a filename or path, use it as public ID
	if s.folder != "" {
		return fmt.Sprintf("%s/%s", s.folder, s.getPublicID(path))
	}
	return s.getPublicID(path)
}
