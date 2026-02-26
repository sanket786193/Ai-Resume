package postgres

import (
	"context"
	"database/sql"

	"resume/internal/domain/entities"

	"github.com/lib/pq"
)

// OfferRepo persists offers.
type OfferRepo struct {
	db *sql.DB
}

// NewOfferRepo returns a new OfferRepo.
func NewOfferRepo(db *sql.DB) *OfferRepo {
	return &OfferRepo{db: db}
}

// Create inserts an offer.
func (r *OfferRepo) Create(ctx context.Context, o *entities.Offer) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO offers (id, ats_id, amount, currency, starts_at, status, responded_at, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		o.ID, o.ATSID, o.Amount, o.Currency, o.StartsAt, o.Status, o.RespondedAt, o.CreatedAt, o.UpdatedAt)
	return err
}

// GetByID returns an offer by ID.
func (r *OfferRepo) GetByID(ctx context.Context, id string) (*entities.Offer, error) {
	var o entities.Offer
	err := r.db.QueryRowContext(ctx,
		`SELECT id, ats_id, amount, currency, starts_at, status, responded_at, created_at, updated_at
		 FROM offers WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&o.ID, &o.ATSID, &o.Amount, &o.Currency, &o.StartsAt, &o.Status, &o.RespondedAt, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// GetByATSID returns the offer for an ATS record (if any).
func (r *OfferRepo) GetByATSID(ctx context.Context, atsID string) (*entities.Offer, error) {
	var o entities.Offer
	err := r.db.QueryRowContext(ctx,
		`SELECT id, ats_id, amount, currency, starts_at, status, responded_at, created_at, updated_at
		 FROM offers WHERE ats_id = $1 AND deleted_at IS NULL`, atsID).
		Scan(&o.ID, &o.ATSID, &o.Amount, &o.Currency, &o.StartsAt, &o.Status, &o.RespondedAt, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// ListByATSIDs returns offers for any of the given ATS record IDs (for HR list).
func (r *OfferRepo) ListByATSIDs(ctx context.Context, atsIDs []string) ([]*entities.Offer, error) {
	if len(atsIDs) == 0 {
		return nil, nil
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, ats_id, amount, currency, starts_at, status, responded_at, created_at, updated_at
		 FROM offers WHERE ats_id = ANY($1) AND deleted_at IS NULL ORDER BY created_at DESC`,
		pq.Array(atsIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*entities.Offer
	for rows.Next() {
		var o entities.Offer
		if err := rows.Scan(&o.ID, &o.ATSID, &o.Amount, &o.Currency, &o.StartsAt, &o.Status, &o.RespondedAt, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &o)
	}
	return list, rows.Err()
}

// UpdateStatus updates offer status and responded_at when accept/reject.
func (r *OfferRepo) UpdateStatus(ctx context.Context, id, status string, respondedAt interface{}) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE offers SET status = $1, responded_at = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 AND deleted_at IS NULL`,
		status, respondedAt, id)
	return err
}

// UpdateStatusTx updates offer status within a transaction.
func (r *OfferRepo) UpdateStatusTx(tx *sql.Tx, ctx context.Context, id, status string, respondedAt interface{}) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE offers SET status = $1, responded_at = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 AND deleted_at IS NULL`,
		status, respondedAt, id)
	return err
}
