package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/google/uuid"
)

// Config for Supabase Storage.
type Config struct {
	URL            string // e.g. https://xxx.supabase.co
	ServiceRoleKey string
	Bucket         string
}

// Client uploads files to Supabase Storage.
type Client struct {
	cfg Config
}

// NewClient creates a Supabase Storage client.
func NewClient(cfg Config) *Client {
	return &Client{cfg: cfg}
}

// UploadResult returned after a successful upload.
type UploadResult struct {
	Key string `json:"Key"`
}

// Upload uploads file content to Supabase Storage and returns the public URL.
func (c *Client) Upload(ctx context.Context, fileContent []byte, fileName string) (*UploadResult, error) {
	url := strings.TrimSuffix(c.cfg.URL, "/")
	if url == "" || c.cfg.ServiceRoleKey == "" || c.cfg.Bucket == "" {
		return nil, fmt.Errorf("supabase storage: missing config")
	}

	// Unique path to avoid 409 Duplicate
	ext := path.Ext(fileName)
	base := strings.TrimSuffix(path.Base(fileName), ext)
	if base == "" {
		base = "resume"
	}
	objectPath := fmt.Sprintf("%s_%s%s", base, uuid.New().String(), ext)

	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", url, c.cfg.Bucket, objectPath)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadURL, bytes.NewReader(fileContent))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.cfg.ServiceRoleKey)
	req.Header.Set("apikey", c.cfg.ServiceRoleKey)
	req.Header.Set("Content-Type", "application/pdf")
	req.Header.Set("x-upsert", "true")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("supabase storage upload: %s %s", resp.Status, string(body))
	}

	var result UploadResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("supabase storage parse response: %w", err)
	}
	return &result, nil
}

// PublicURL returns the public URL for an object key (Key from UploadResult).
// Key may be "path" or "bucket/path"; both are supported.
func (c *Client) PublicURL(key string) string {
	url := strings.TrimSuffix(c.cfg.URL, "/")
	if key != "" && !strings.HasPrefix(key, c.cfg.Bucket+"/") {
		key = c.cfg.Bucket + "/" + key
	}
	return fmt.Sprintf("%s/storage/v1/object/public/%s", url, key)
}
