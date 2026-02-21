package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
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

// sanitizeObjectBase normalizes a filename base for storage paths so that spaces
// and parentheses do not break signed/public URLs. Example:
// "Golang_resume_latest (1)" -> "golang_resume_latest_1"
var multiUnderscore = regexp.MustCompile(`_+`)

func sanitizeObjectBase(base string) string {
	base = strings.TrimSpace(base)
	if base == "" {
		return "resume"
	}
	// Replace spaces and parentheses with underscores; remove other problematic chars
	base = strings.ReplaceAll(base, " ", "_")
	base = strings.ReplaceAll(base, "(", "_")
	base = strings.ReplaceAll(base, ")", "")
	base = multiUnderscore.ReplaceAllString(base, "_")
	base = strings.Trim(base, "_")
	if base == "" {
		return "resume"
	}
	return strings.ToLower(base)
}

// Upload uploads file content to Supabase Storage and returns the public URL.
func (c *Client) Upload(ctx context.Context, fileContent []byte, fileName string) (*UploadResult, error) {
	baseURL := strings.TrimSuffix(c.cfg.URL, "/")
	if baseURL == "" || c.cfg.ServiceRoleKey == "" || c.cfg.Bucket == "" {
		return nil, fmt.Errorf("supabase storage: missing config")
	}

	// Unique path to avoid 409 Duplicate; sanitize base so path has no spaces/parentheses
	ext := strings.ToLower(path.Ext(fileName))
	if ext == "" {
		ext = ".pdf"
	}
	rawBase := strings.TrimSuffix(path.Base(fileName), path.Ext(fileName))
	base := sanitizeObjectBase(rawBase)
	objectPath := fmt.Sprintf("%s_%s%s", base, uuid.New().String(), ext)

	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", baseURL, c.cfg.Bucket, objectPath)
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
// Path is case-sensitive and must match storage exactly: /storage/v1/object/public/{bucket}/{path}.
func (c *Client) PublicURL(key string) string {
	baseURL := strings.TrimSuffix(c.cfg.URL, "/")
	if key == "" {
		return ""
	}
	// Key from Supabase is often just the object path; ensure full bucket/path for URL
	if !strings.HasPrefix(key, c.cfg.Bucket+"/") {
		key = c.cfg.Bucket + "/" + key
	}
	// Encode each path segment (not the whole path) so slashes stay as slashes
	parts := strings.Split(key, "/")
	for i, p := range parts {
		parts[i] = url.PathEscape(p)
	}
	encodedPath := strings.Join(parts, "/")
	return fmt.Sprintf("%s/storage/v1/object/public/%s", baseURL, encodedPath)
}
