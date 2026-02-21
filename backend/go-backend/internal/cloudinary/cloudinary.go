package cloudinary

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// Config for Cloudinary.
type Config struct {
	CloudName  string
	APIKey     string
	APISecret  string
	UploadFolder string // e.g. "resumes"
}

// UploadResult returned after a successful upload.
type UploadResult struct {
	SecureURL string `json:"secure_url"`
	PublicID  string `json:"public_id"`
}

// Client uploads files to Cloudinary.
type Client struct {
	cfg Config
}

// NewClient creates a Cloudinary client.
func NewClient(cfg Config) *Client {
	return &Client{cfg: cfg}
}

// Upload uploads file content to Cloudinary (raw/unsigned would need upload_preset;
// we use signed upload). Returns the secure URL and public_id.
func (c *Client) Upload(ctx context.Context, fileContent []byte, fileName string) (*UploadResult, error) {
	if c.cfg.CloudName == "" || c.cfg.APIKey == "" || c.cfg.APISecret == "" {
		return nil, fmt.Errorf("cloudinary: missing config")
	}

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	// Add file
	part, err := w.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, bytes.NewReader(fileContent)); err != nil {
		return nil, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	params := map[string]string{
		"timestamp": timestamp,
		"api_key":   c.cfg.APIKey,
	}
	if c.cfg.UploadFolder != "" {
		params["folder"] = c.cfg.UploadFolder
	}

	signature := c.signParams(params)
	_ = w.WriteField("timestamp", timestamp)
	_ = w.WriteField("api_key", c.cfg.APIKey)
	_ = w.WriteField("signature", signature)
	if c.cfg.UploadFolder != "" {
		_ = w.WriteField("folder", c.cfg.UploadFolder)
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/raw/upload", c.cfg.CloudName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("cloudinary upload: %s %s", resp.Status, string(b))
	}

	// Parse JSON response: { "secure_url": "...", "public_id": "..." }
	var result struct {
		SecureURL string `json:"secure_url"`
		PublicID  string `json:"public_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &UploadResult{SecureURL: result.SecureURL, PublicID: result.PublicID}, nil
}

// signParams builds Cloudinary signature: sha1(sorted params string + api_secret).
func (c *Client) signParams(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var s string
	for i, k := range keys {
		if i > 0 {
			s += "&"
		}
		s += k + "=" + params[k]
	}
	h := sha1.Sum([]byte(s + c.cfg.APISecret))
	return hex.EncodeToString(h[:])
}
