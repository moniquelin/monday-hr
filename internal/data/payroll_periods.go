package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

var (
	ErrPayrollPeriodOverlap   = errors.New("overlapping date with existing period")
	ErrPayrollPeriodDateOrder = errors.New("start date is greater than end date")
)

// Payroll struct represents attendance data of one date
type PayrollPeriod struct {
	ID        int64  `json:"id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Status    string `json:"status"`

	ProcessedAt time.Time `json:"processed_at"`
	ProcessedBy time.Time `json:"processed_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int64     `json:"created_by"`
	UpdatedBy int64     `json:"updated_by"`
}

type PayrollPeriodModel struct {
	DB *sql.DB
}

// CheckOverlap checks whether a payroll period overlaps with any existing period.
func (m PayrollPeriodModel) CheckOverlap(startDate, endDate string) error {
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
func (m PayrollPeriodModel) Insert(p PayrollPeriod) error {
	// Check if new dates overlap with existing periods
	if err := m.CheckOverlap(p.StartDate, p.EndDate); err != nil {
		return err
	}

	query := `
    INSERT INTO payroll_periods (start_date, end_date, created_by, updated_by)
    VALUES ($1, $2, $3, $4)
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Insert new payroll period
	_, err := m.DB.ExecContext(ctx, query, p.StartDate, p.EndDate, p.CreatedBy, p.UpdatedBy)
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
