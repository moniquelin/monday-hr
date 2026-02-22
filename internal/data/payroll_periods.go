package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

var (
	ErrPayrollPeriodOverlap    = errors.New("overlapping date with existing period")
	ErrPayrollPeriodDateOrder  = errors.New("start date is greater than end date")
	ErrPayrollAlreadyProcessed = errors.New("payroll preiod has already been processed")
)

// PayrollPeriod struct represents a payroll period
type PayrollPeriod struct {
	ID        int64 `json:"id"`
	StartDate time.Time
	EndDate   time.Time
	Status    string `json:"status"`

	ProcessedAt *time.Time `json:"processed_at"`
	ProcessedBy *int64     `json:"processed_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int64     `json:"created_by"`
	UpdatedBy int64     `json:"updated_by"`
}

type PayrollPeriodModel struct {
	DB *sql.DB
}

// CheckOverlap checks whether a payroll period overlaps with any existing period.
func (m PayrollPeriodModel) CheckOverlap(startDate, endDate time.Time) error {
	query := `
	SELECT EXISTS (
		SELECT 1
		FROM payroll_periods
		WHERE start_date <= $2
		AND end_date >= $1
	)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var overlap bool
	err := m.DB.QueryRowContext(ctx, query, startDate, endDate).
		Scan(&overlap)
	if err != nil {
		return err
	}

	// If the table already contains a record
	if overlap {
		return ErrPayrollPeriodOverlap
	}
	return nil
}

// Insert new payroll period in the database
func (m PayrollPeriodModel) Insert(startDate, endDate time.Time, userId int64) error {
	// Check if new dates overlap with existing period
	if err := m.CheckOverlap(startDate, endDate); err != nil {
		return err
	}

	query := `
    INSERT INTO payroll_periods (start_date, end_date, created_by, updated_by)
    VALUES ($1, $2, $3, $3)
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Insert new payroll period
	_, err := m.DB.ExecContext(ctx, query, startDate, endDate, userId)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Constraint {
			case "payroll_periods_prevent_date_overlap":
				return ErrPayrollPeriodOverlap
			case "chk_period_date_order":
				return ErrPayrollPeriodDateOrder
			}
		}
		return err
	}
	return nil
}

// Get payroll period by ID from the database
func (m PayrollPeriodModel) Get(id int64) (*PayrollPeriod, error) {
	query := `
		SELECT id, start_date, end_date, status,
		       processed_at, processed_by,
		       created_at, created_by, updated_at, updated_by
		FROM payroll_periods
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var p PayrollPeriod

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.StartDate,
		&p.EndDate,
		&p.Status,
		&p.ProcessedAt,
		&p.ProcessedBy,
		&p.CreatedAt,
		&p.CreatedBy,
		&p.UpdatedAt,
		&p.UpdatedBy,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &p, nil
}

// MarkAsProcessed sets payroll period status to processed.
// Caller must ensure the period exists and is still draft.
func (m PayrollPeriodModel) MarkAsProcessed(id int64, processedBy int64) error {
	query := `
		UPDATE payroll_periods
		SET status = 'processed',
		    processed_at = NOW(),
		    processed_by = $2,
			updated_at = NOW(),
		    updated_by = $2
		WHERE id = $1
			AND status = 'draft'
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id, processedBy)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrRecordNotFound
	}

	return nil
}
