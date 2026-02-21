package notifications

import (
	"context"
	"log"
)

// Service handles async notifications (email, etc.); stub for production integration.
type Service struct{}

// NewService creates a notifications service.
func NewService() *Service {
	return &Service{}
}

// SendApplicationReceived notifies candidate that application was received (async).
func (s *Service) SendApplicationReceived(ctx context.Context, candidateEmail, jobTitle string) {
	log.Printf("[notifications] application received: %s for job %s", candidateEmail, jobTitle)
	// TODO: enqueue email or call email provider
}

// SendOfferCreated notifies candidate that an offer was created.
func (s *Service) SendOfferCreated(ctx context.Context, candidateEmail, jobTitle string) {
	log.Printf("[notifications] offer created: %s for job %s", candidateEmail, jobTitle)
}

// SendOfferAccepted notifies HR that offer was accepted.
func (s *Service) SendOfferAccepted(ctx context.Context, hrEmail, candidateName string) {
	log.Printf("[notifications] offer accepted: candidate %s", candidateName)
}
