package notifications

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"resume/internal/config"
)

// Service handles async notifications (email via SMTP when enabled).
type Service struct {
	smtp *config.SMTPConfig
}

// NewService creates a notifications service (log-only).
func NewService() *Service {
	return &Service{}
}

// NewServiceWithSMTP creates a notifications service that sends email when smtp.Enabled is true.
func NewServiceWithSMTP(smtp *config.SMTPConfig) *Service {
	return &Service{smtp: smtp}
}

// sendEmail sends a plain-text email via SMTP. No-op if SMTP not configured or disabled.
func (s *Service) sendEmail(to, subject, body string) {
	if s.smtp == nil || !s.smtp.Enabled || s.smtp.Host == "" || s.smtp.Email == "" || s.smtp.Password == "" {
		log.Printf("[notifications] email (no SMTP): to=%s subject=%s", to, subject)
		return
	}
	addr := s.smtp.Host + ":" + s.smtp.Port
	auth := smtp.PlainAuth("", s.smtp.Email, strings.TrimSpace(s.smtp.Password), s.smtp.Host)
	from := s.smtp.Email
	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" + body + "\r\n")
	if err := smtp.SendMail(addr, auth, from, []string{to}, msg); err != nil {
		log.Printf("[notifications] SMTP send failed to %s: %v", to, err)
		return
	}
	log.Printf("[notifications] email sent to %s: %s", to, subject)
}

// SendApplicationReceived notifies candidate that application was received (async).
func (s *Service) SendApplicationReceived(ctx context.Context, candidateEmail, jobTitle string) {
	log.Printf("[notifications] application received: %s for job %s", candidateEmail, jobTitle)
	s.sendEmail(candidateEmail, "Application received – "+jobTitle,
		fmt.Sprintf("Hi,\n\nWe have received your application for the position: %s.\n\nWe will review it and get back to you.\n\nBest regards,\nRecruitment Team", jobTitle))
}

// SendShortlisted notifies candidate they were shortlisted (Phase 3).
func (s *Service) SendShortlisted(ctx context.Context, candidateEmail, jobTitle string) {
	log.Printf("[notifications] shortlisted: %s for job %s", candidateEmail, jobTitle)
	s.sendEmail(candidateEmail, "You have been shortlisted – "+jobTitle,
		fmt.Sprintf("Hi,\n\nCongratulations! You have been shortlisted for the position: %s.\n\nWe will contact you soon regarding the next steps.\n\nBest regards,\nRecruitment Team", jobTitle))
}

// SendRejected notifies candidate their application was rejected (Phase 3).
func (s *Service) SendRejected(ctx context.Context, candidateEmail, jobTitle string) {
	log.Printf("[notifications] rejected: %s for job %s", candidateEmail, jobTitle)
	s.sendEmail(candidateEmail, "Update on your application – "+jobTitle,
		fmt.Sprintf("Hi,\n\nThank you for your interest in the position: %s.\n\nAfter careful consideration, we have decided to move forward with other candidates at this time. We encourage you to apply for future openings.\n\nBest regards,\nRecruitment Team", jobTitle))
}

// SendInterviewScheduled notifies candidate that an interview was scheduled (Phase 3).
func (s *Service) SendInterviewScheduled(ctx context.Context, candidateEmail, jobTitle string) {
	log.Printf("[notifications] interview scheduled: %s for job %s", candidateEmail, jobTitle)
	s.sendEmail(candidateEmail, "Interview scheduled – "+jobTitle,
		fmt.Sprintf("Hi,\n\nAn interview has been scheduled for the position: %s.\n\nYou will receive further details (date, time, and link/location) separately. Please confirm your availability.\n\nBest regards,\nRecruitment Team", jobTitle))
}

// SendOfferLetterReady notifies candidate that offer letter is ready to download (Phase 3).
func (s *Service) SendOfferLetterReady(ctx context.Context, candidateEmail, jobTitle string) {
	log.Printf("[notifications] offer letter ready: %s for job %s", candidateEmail, jobTitle)
	s.sendEmail(candidateEmail, "Offer letter ready – "+jobTitle,
		fmt.Sprintf("Hi,\n\nCongratulations! Your offer letter for the position %s is ready.\n\nPlease log in to the applicant portal to view and download your offer letter.\n\nBest regards,\nRecruitment Team", jobTitle))
}

// SendOfferCreated notifies candidate that an offer was created.
func (s *Service) SendOfferCreated(ctx context.Context, candidateEmail, jobTitle string) {
	log.Printf("[notifications] offer created: %s for job %s", candidateEmail, jobTitle)
	s.sendEmail(candidateEmail, "Offer – "+jobTitle,
		fmt.Sprintf("Hi,\n\nAn offer has been created for the position: %s.\n\nPlease log in to the applicant portal to view and respond to the offer.\n\nBest regards,\nRecruitment Team", jobTitle))
}

// SendOfferAccepted notifies HR that offer was accepted.
func (s *Service) SendOfferAccepted(ctx context.Context, hrEmail, candidateName string) {
	log.Printf("[notifications] offer accepted: candidate %s", candidateName)
	s.sendEmail(hrEmail, "Offer accepted – "+candidateName,
		fmt.Sprintf("Hi,\n\n%s has accepted the offer.\n\nBest regards,\nATS", candidateName))
}
