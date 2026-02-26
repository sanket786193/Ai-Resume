package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client calls the Python AI service (Ollama/ADK) over HTTP.
type Client struct {
	baseURL    string
	httpClient *http.Client
	enabled    bool
}

// JobRequirements is optional context for matching (Phase 2).
type JobRequirements struct {
	Skills          []string `json:"skills,omitempty"`
	ExperienceLevel string   `json:"experience_level,omitempty"`
	Qualification   string   `json:"qualification,omitempty"`
}

// ScreenRequest is sent to the AI service.
type ScreenRequest struct {
	ResumePathOrContent string           `json:"resume_path_or_content"`
	JobDescription     string           `json:"job_description"`
	JobRequirements    *JobRequirements `json:"job_requirements,omitempty"`
}

// ScreenResponse is returned by the AI service (full ATS evaluation).
type ScreenResponse struct {
	SkillMatchScore     float64  `json:"skill_match_score"`
	RankingScore        float64  `json:"ranking_score"`
	Qualified           bool     `json:"qualified"`
	ATSScore            *int     `json:"ats_score,omitempty"`
	SkillMatchPct       *int     `json:"skill_match_pct,omitempty"`
	MissingSkills       []string `json:"missing_skills,omitempty"`
	ExperienceMatch     *string  `json:"experience_match,omitempty"`
	ExperienceWarnings  []string `json:"experience_warnings,omitempty"`
	KeywordMatches      []string `json:"keyword_matches,omitempty"`
	SemanticMatches     []string `json:"semantic_matches,omitempty"`
	Summary             *string  `json:"summary,omitempty"`
	ModelVersion        *string  `json:"model_version,omitempty"`
}

// New creates an AI client.
func New(baseURL string, timeoutSec int, enabled bool) *Client {
	if timeoutSec <= 0 {
		timeoutSec = 60
	}
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
		},
		enabled: enabled,
	}
}

// ParseRequest is sent to POST /parse.
type ParseRequest struct {
	ResumePathOrContent string `json:"resume_path_or_content"`
}

// ParseResponse is returned by the AI service parse endpoint.
type ParseResponse struct {
	RawText     string          `json:"raw_text"`
	Parsed      json.RawMessage `json:"parsed"`
	CleanedText string          `json:"cleaned_text"`
}

// EmbedRequest is sent to POST /embed.
type EmbedRequest struct {
	Text string `json:"text"`
}

// EmbedResponse is returned by the AI service embed endpoint.
type EmbedResponse struct {
	Embedding   []float64 `json:"embedding"`
	ModelVersion *string  `json:"model_version,omitempty"`
}

// Embed calls the Python service to get embedding vector for text (e.g. cleaned resume).
func (c *Client) Embed(ctx context.Context, text string) ([]float64, string, error) {
	if !c.enabled || text == "" {
		return nil, "", nil
	}
	reqBody := EmbedRequest{Text: text}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/embed", bytes.NewReader(body))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("ai embed returned %d", resp.StatusCode)
	}
	var out EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, "", err
	}
	modelVer := ""
	if out.ModelVersion != nil {
		modelVer = *out.ModelVersion
	}
	return out.Embedding, modelVer, nil
}

// Parse calls the Python service to parse resume (URL or raw text); returns raw_text, parsed_json, cleaned_text.
func (c *Client) Parse(ctx context.Context, resumePathOrContent string) (*ParseResponse, error) {
	if !c.enabled {
		return &ParseResponse{}, nil
	}
	reqBody := ParseRequest{ResumePathOrContent: resumePathOrContent}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/parse", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ai parse returned %d", resp.StatusCode)
	}
	var out ParseResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ScreenResume calls the Python service for resume screening; implements ats.AIClient.
// jobRequirements is optional; when set, the AI service uses it for skill/experience/qualification matching.
func (c *Client) ScreenResume(ctx context.Context, resumeContentOrPath, jobDescription string, jobRequirements *JobRequirements) (*ScreenResponse, error) {
	if !c.enabled {
		return &ScreenResponse{}, nil
	}
	reqBody := ScreenRequest{
		ResumePathOrContent: resumeContentOrPath,
		JobDescription:     jobDescription,
		JobRequirements:    jobRequirements,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/screen", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ai service returned %d", resp.StatusCode)
	}
	var out ScreenResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
