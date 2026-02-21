package offers

import (
	"context"
	"database/sql"
	"time"

	"resume/internal/domain/entities"
	domainerrors "resume/internal/domain/errors"
	"resume/internal/domain/enums"
	"resume/internal/storage/postgres"
	"resume/internal/storage/tx"

	"github.com/google/uuid"
)

// Service contains offer creation and accept/reject logic; ATS → HIRED only after accept.
type Service struct {
	offerRepo *postgres.OfferRepo
	atsRepo   *postgres.ATSRepo
	txRunner  *tx.Runner
}

// NewService creates an offers service.
func NewService(offerRepo *postgres.OfferRepo, atsRepo *postgres.ATSRepo, txRunner *tx.Runner) *Service {
	return &Service{offerRepo: offerRepo, atsRepo: atsRepo, txRunner: txRunner}
}

// Initiate creates an offer for an ATS record (HR).
func (s *Service) Initiate(ctx context.Context, atsID, amount, currency string, startsAt *time.Time) (*entities.Offer, error) {
	rec, err := s.atsRepo.GetByID(ctx, atsID)
	if err != nil || rec == nil {
		return nil, &domainerrors.NotFoundError{Resource: "ats_record", ID: atsID}
	}
	if rec.Status != enums.ATSShortlisted && rec.Status != enums.ATSInterview {
		return nil, &domainerrors.ValidationError{Field: "ats_id", Message: "candidate must be shortlisted or in interview to receive offer"}
	}
	existing, _ := s.offerRepo.GetByATSID(ctx, atsID)
	if existing != nil {
		return nil, &domainerrors.ConflictError{Message: "offer already exists for this application"}
	}
	now := time.Now()
	offer := &entities.Offer{
		ID:        uuid.New().String(),
		ATSID:     atsID,
		Amount:    amount,
		Currency:  currency,
		StartsAt:  startsAt,
		Status:    "PENDING",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if offer.Currency == "" {
		offer.Currency = "USD"
	}
	if err := s.offerRepo.Create(ctx, offer); err != nil {
		return nil, err
	}
	return offer, nil
}

// Accept sets offer to ACCEPTED and ATS to HIRED (atomic transaction).
func (s *Service) Accept(ctx context.Context, offerID string) error {
	offer, err := s.offerRepo.GetByID(ctx, offerID)
	if err != nil || offer == nil {
		return &domainerrors.NotFoundError{Resource: "offer", ID: offerID}
	}
	if offer.Status != "PENDING" {
		return &domainerrors.ConflictError{Message: "offer already responded"}
	}
	now := time.Now()
	return s.txRunner.Run(ctx, func(tx *sql.Tx) error {
		if err := s.offerRepo.UpdateStatusTx(tx, ctx, offerID, "ACCEPTED", now); err != nil {
			return err
		}
		return s.atsRepo.UpdateStatusTx(tx, ctx, offer.ATSID, enums.ATSHired)
	})
}

// Reject sets offer to REJECTED.
func (s *Service) Reject(ctx context.Context, offerID string) error {
	offer, err := s.offerRepo.GetByID(ctx, offerID)
	if err != nil || offer == nil {
		return &domainerrors.NotFoundError{Resource: "offer", ID: offerID}
	}
	if offer.Status != "PENDING" {
		return &domainerrors.ConflictError{Message: "offer already responded"}
	}
	now := time.Now()
	return s.offerRepo.UpdateStatus(ctx, offerID, "REJECTED", now)
}

// GetByID returns an offer by ID.
func (s *Service) GetByID(ctx context.Context, id string) (*entities.Offer, error) {
	o, err := s.offerRepo.GetByID(ctx, id)
	if err != nil || o == nil {
		return nil, &domainerrors.NotFoundError{Resource: "offer", ID: id}
	}
	return o, nil
}

// GetByATSID returns the offer for an ATS record.
func (s *Service) GetByATSID(ctx context.Context, atsID string) (*entities.Offer, error) {
	return s.offerRepo.GetByATSID(ctx, atsID)
}
