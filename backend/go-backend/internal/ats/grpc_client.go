package ats

import (
	"context"

	atspb "resume/proto/ats"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// grpcAIClient implements AIClient using gRPC (Python AI service).
type grpcAIClient struct {
	conn   *grpc.ClientConn
	client atspb.AIScreeningServiceClient
}

// NewGRPCAIClient dials the AI service gRPC server and returns an AIClient.
func NewGRPCAIClient(target string, timeoutSec int, enabled bool) (AIClient, error) {
	if !enabled {
		return &noopAIClient{}, nil
	}
	if timeoutSec <= 0 {
		timeoutSec = 60
	}
	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return &grpcAIClient{
		conn:   conn,
		client: atspb.NewAIScreeningServiceClient(conn),
	}, nil
}

// ScreenResume implements AIClient.
func (c *grpcAIClient) ScreenResume(ctx context.Context, resumeContentOrPath, jobDescription string) (*AIScreenResult, error) {
	if c.client == nil {
		return &AIScreenResult{}, nil
	}
	req := &atspb.ScreenRequest{
		ResumePathOrContent: resumeContentOrPath,
		JobDescription:      jobDescription,
	}
	resp, err := c.client.Screen(ctx, req)
	if err != nil {
		return nil, err
	}
	return &AIScreenResult{
		SkillMatchScore: resp.GetSkillMatchScore(),
		RankingScore:    resp.GetRankingScore(),
		Qualified:       resp.GetQualified(),
	}, nil
}

// Parse implements AIClient (no-op for gRPC; use HTTP adapter for parse + persist).
func (c *grpcAIClient) Parse(ctx context.Context, resumePathOrContent string) (rawText string, parsedJSON []byte, cleanedText string, err error) {
	return "", nil, "", nil
}

// noopAIClient returns empty result when AI is disabled.
type noopAIClient struct{}

func (noopAIClient) ScreenResume(context.Context, string, string) (*AIScreenResult, error) {
	return &AIScreenResult{}, nil
}

func (noopAIClient) Parse(context.Context, string) (string, []byte, string, error) {
	return "", nil, "", nil
}
