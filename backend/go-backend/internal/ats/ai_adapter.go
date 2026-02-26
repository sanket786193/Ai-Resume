package ats

import (
	"context"

	"resume/internal/ai/client"
)

// aiAdapter adapts ai/client.Client to AIClient interface.
type aiAdapter struct {
	client *client.Client
}

// ScreenResume calls the AI service and converts response to AIScreenResult.
func (a *aiAdapter) ScreenResume(ctx context.Context, resumeContentOrPath, jobDescription string, jobRequirements *client.JobRequirements) (*AIScreenResult, error) {
	resp, err := a.client.ScreenResume(ctx, resumeContentOrPath, jobDescription, jobRequirements)
	if err != nil {
		return nil, err
	}
	out := &AIScreenResult{
		SkillMatchScore: resp.SkillMatchScore,
		RankingScore:    resp.RankingScore,
		Qualified:       resp.Qualified,
	}
	if resp.ATSScore != nil {
		out.ATSScore = resp.ATSScore
	}
	if resp.SkillMatchPct != nil {
		out.SkillMatchPct = resp.SkillMatchPct
	}
	if len(resp.MissingSkills) > 0 {
		out.MissingSkills = resp.MissingSkills
	}
	if resp.ExperienceMatch != nil {
		out.ExperienceMatch = resp.ExperienceMatch
	}
	if len(resp.ExperienceWarnings) > 0 {
		out.ExperienceWarnings = resp.ExperienceWarnings
	}
	if len(resp.KeywordMatches) > 0 {
		out.KeywordMatches = resp.KeywordMatches
	}
	if len(resp.SemanticMatches) > 0 {
		out.SemanticMatches = resp.SemanticMatches
	}
	if resp.Summary != nil {
		out.Summary = resp.Summary
	}
	if resp.ModelVersion != nil {
		out.ModelVersion = resp.ModelVersion
	}
	return out, nil
}

// Parse calls the AI service parse endpoint and returns raw text, parsed JSON, and cleaned text.
func (a *aiAdapter) Parse(ctx context.Context, resumePathOrContent string) (rawText string, parsedJSON []byte, cleanedText string, err error) {
	resp, err := a.client.Parse(ctx, resumePathOrContent)
	if err != nil {
		return "", nil, "", err
	}
	return resp.RawText, resp.Parsed, resp.CleanedText, nil
}

// NewAIAdapter returns an AIClient that uses the given HTTP client.
func NewAIAdapter(c *client.Client) AIClient {
	return &aiAdapter{client: c}
}
